package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	apiErrors "github.com/Archisman-Mridha/x-clone/backend/microservices/users/internal/errors"
	"github.com/Archisman-Mridha/x-clone/backend/microservices/users/internal/postgresql/generated"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/connectors"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

type (
	UsersRepository interface {
		Create(ctx context.Context, args *CreateUserArgs) (int32, error)

		FindByEmail(ctx context.Context, email string) (*FindUserByOperationOutput, error)
		FindByUsername(ctx context.Context, username string) (*FindUserByOperationOutput, error)

		IDExists(ctx context.Context, id int32) (bool, error)
	}

	CreateUserArgs struct {
		Name,
		Email,
		Username,
		HashedPassword string
	}

	FindUserByOperationOutput struct {
		ID             int32
		HashedPassword string
	}
)

type UsersPostgresRepository struct {
	*connectors.PostgresConnector
	queries *generated.Queries
}

func NewUsersPostgresRepository(postgresConnector *connectors.PostgresConnector) UsersRepository {
	queries := generated.New(postgresConnector.GetConnection())

	return &UsersPostgresRepository{
		postgresConnector,
		queries,
	}
}

func (u *UsersPostgresRepository) Create(ctx context.Context,
	args *CreateUserArgs,
) (int32, error) {
	userID, err := u.queries.CreateUser(ctx, (*generated.CreateUserParams)(args))
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && (pgErr.Code == pgerrcode.UniqueViolation) {
			switch pgErr.ColumnName {
			case "email":
				return 0, apiErrors.ErrDuplicateEmail

			case "username":
				return 0, apiErrors.ErrDuplicateUsername
			}
		}

		return 0, utils.WrapError(err)
	}
	return userID, nil
}

func (u *UsersPostgresRepository) FindByEmail(ctx context.Context,
	email string,
) (*FindUserByOperationOutput, error) {
	userDetails, err := u.queries.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apiErrors.ErrUserNotFound
		}

		return nil, utils.WrapError(err)
	}
	return (*FindUserByOperationOutput)(userDetails), nil
}

func (u *UsersPostgresRepository) FindByUsername(ctx context.Context,
	username string,
) (*FindUserByOperationOutput, error) {
	userDetails, err := u.queries.FindUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apiErrors.ErrUserNotFound
		}

		return nil, utils.WrapError(err)
	}
	return (*FindUserByOperationOutput)(userDetails), nil
}

func (u *UsersPostgresRepository) IDExists(ctx context.Context, id int32) (bool, error) {
	_, err := u.queries.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, apiErrors.ErrUserNotFound
		}

		return false, utils.WrapError(err)
	}
	return true, nil
}
