package followships

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Archisman-Mridha/x-clone/backend/microservices/followships/internal/postgresql/generated"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/connectors"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

type (
	FollowshipsRepository interface {
		CreateFollowship(ctx context.Context, args *FollowshipOperationArgs) error
		DeleteFollowship(ctx context.Context, args *FollowshipOperationArgs) error

		FollowshipExists(ctx context.Context, args *FollowshipOperationArgs) (bool, error)

		GetFollowers(ctx context.Context, args *GetFollowersArgs) ([]int32, error)
		GetFollowees(ctx context.Context, args *GetFolloweesArgs) ([]int32, error)
		GetFollowerAndFolloweeCounts(ctx context.Context,
			userID int32,
		) (*GetFollowerAndFolloweeCountsOutput, error)
	}

	Followship struct {
		ID,

		FollowerID,
		FolloweeID int32
	}
)

type FollowshipsPostgresRepository struct {
	*connectors.PostgresConnector
	queries *generated.Queries
}

func NewFollowshipsPostgresRepository(
	postgresConnector *connectors.PostgresConnector,
) FollowshipsRepository {
	queries := generated.New(postgresConnector.GetConnection())

	return &FollowshipsPostgresRepository{
		postgresConnector,
		queries,
	}
}

func (f *FollowshipsPostgresRepository) CreateFollowship(ctx context.Context,
	args *FollowshipOperationArgs,
) error {
	err := f.queries.CreateFollowship(ctx, (*generated.CreateFollowshipParams)(args))
	if err != nil {
		return utils.WrapError(err)
	}
	return nil
}

func (f *FollowshipsPostgresRepository) DeleteFollowship(ctx context.Context,
	args *FollowshipOperationArgs,
) error {
	err := f.queries.DeleteFollowship(ctx, (*generated.DeleteFollowshipParams)(args))
	if err != nil {
		return utils.WrapError(err)
	}
	return nil
}

func (f *FollowshipsPostgresRepository) FollowshipExists(ctx context.Context,
	args *FollowshipOperationArgs,
) (bool, error) {
	err := f.queries.GetFollowship(ctx, (*generated.GetFollowshipParams)(args))
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && (pgErr.Code == pgerrcode.NoDataFound) {
			return false, nil
		}

		return false, utils.WrapError(err)
	}
	return true, nil
}

func (f *FollowshipsPostgresRepository) GetFollowers(ctx context.Context,
	args *GetFollowersArgs,
) ([]int32, error) {
	followers, err := f.queries.GetFollowers(ctx, &generated.GetFollowersParams{
		FolloweeID: args.FolloweeID,

		Offset: int32(args.PageRequest.Offset),
		Limit:  int32(args.PageRequest.Limit),
	})
	if err != nil {
		return nil, utils.WrapError(err)
	}
	return followers, nil
}

func (f *FollowshipsPostgresRepository) GetFollowees(ctx context.Context,
	args *GetFolloweesArgs,
) ([]int32, error) {
	followees, err := f.queries.GetFollowees(ctx, &generated.GetFolloweesParams{
		FollowerID: args.FollowerID,

		Offset: int32(args.PageRequest.Offset),
		Limit:  int32(args.PageRequest.Limit),
	})
	if err != nil {
		return nil, utils.WrapError(err)
	}
	return followees, nil
}

func (f *FollowshipsPostgresRepository) GetFollowerAndFolloweeCounts(ctx context.Context,
	userID int32,
) (*GetFollowerAndFolloweeCountsOutput, error) {
	followerAndFolloweeCounts, err := f.queries.GetFollowerAndFolloweeCounts(ctx, userID)
	if err != nil {
		return nil, utils.WrapError(err)
	}
	return (*GetFollowerAndFolloweeCountsOutput)(followerAndFolloweeCounts), nil
}
