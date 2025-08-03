package profiles

import (
	"context"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Archisman-Mridha/x-clone/backend/microservices/profiles/internal/postgresql/generated"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/connectors"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

type (
	ProfilesRepository interface {
		// This operation needs to be idempotent, since it gets invoked by a DB event processor.
		Create(ctx context.Context, args *CreateProfileArgs) error

		GetPreviews(ctx context.Context, ids []int32) ([]*ProfilePreview, error)
	}

	ProfilePreview struct {
		ID int32
		Name,
		Username string
	}

	Profile struct {
		ProfilePreview
	}
)

type ProfilesPostgresRepository struct {
	*connectors.PostgresConnector
	queries *generated.Queries
}

func NewProfilesPostgresRepository(
	postgresConnector *connectors.PostgresConnector,
) ProfilesRepository {
	queries := generated.New(postgresConnector.GetConnection())

	return &ProfilesPostgresRepository{
		postgresConnector,
		queries,
	}
}

func (p *ProfilesPostgresRepository) Create(ctx context.Context, args *CreateProfileArgs) error {
	err := p.queries.CreateProfile(ctx, (*generated.CreateProfileParams)(args))
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)

		// We'll not error out, if a profile with the given ID already exists.
		// This makes the operation idempotent.
		if ok && (pgErr.Code == pgerrcode.UniqueViolation) && (pgErr.ColumnName == "id") {
			slog.WarnContext(
				ctx,
				"Can't create profile, since it already exists. Most probably duplicate processed a Kafka record.",
				slog.Int("profile-id", int(args.ID)),
			)
			return nil
		}

		return utils.WrapError(err)
	}
	return nil
}

func (p *ProfilesPostgresRepository) GetPreviews(ctx context.Context,
	ids []int32,
) ([]*ProfilePreview, error) {
	profilePreviews := []*ProfilePreview{}

	rows, err := p.queries.GetProfilePreviews(ctx, ids)
	if err != nil {
		return profilePreviews, utils.WrapError(err)
	}

	for _, row := range rows {
		profilePreviews = append(profilePreviews, (*ProfilePreview)(row))
	}
	return profilePreviews, nil
}
