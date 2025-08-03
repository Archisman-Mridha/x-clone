package posts

import (
	"context"
	"time"

	"github.com/Archisman-Mridha/x-clone/backend/microservices/posts/internal/postgresql/generated"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/connectors"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

type (
	PostsRepository interface {
		Create(ctx context.Context, args *CreatePostArgs) (int32, error)

		GetPosts(ctx context.Context, ids []int32) ([]*Post, error)
		GetPostsOfUser(ctx context.Context, args *GetPostsOfUserArgs) ([]*Post, error)
	}

	Post struct {
		ID,
		OwnerID int32

		Description string

		CreatedAt time.Time
	}
)

type PostsPostgresRepository struct {
	*connectors.PostgresConnector
	queries *generated.Queries
}

func NewPostsPostgresRepository(postgresConnector *connectors.PostgresConnector) PostsRepository {
	queries := generated.New(postgresConnector.GetConnection())

	return &PostsPostgresRepository{
		postgresConnector,
		queries,
	}
}

func (p *PostsPostgresRepository) Create(ctx context.Context, args *CreatePostArgs) (int32, error) {
	postID, err := p.queries.CreatePost(ctx, (*generated.CreatePostParams)(args))
	if err != nil {
		return 0, utils.WrapError(err)
	}
	return postID, nil
}

func (p *PostsPostgresRepository) GetPosts(ctx context.Context, ids []int32) ([]*Post, error) {
	posts := []*Post{}

	rows, err := p.queries.GetPosts(ctx, ids)
	if err != nil {
		return posts, utils.WrapError(err)
	}

	for _, row := range rows {
		posts = append(posts, (*Post)(row))
	}
	return posts, nil
}

func (p *PostsPostgresRepository) GetPostsOfUser(ctx context.Context,
	args *GetPostsOfUserArgs,
) ([]*Post, error) {
	posts := []*Post{}

	rows, err := p.queries.GetPostsOfUser(ctx, &generated.GetPostsOfUserParams{
		OwnerID: args.OwnerID,

		Offset: int32(args.PageRequest.Offset),
		Limit:  int32(args.PageRequest.Limit),
	})
	if err != nil {
		return posts, utils.WrapError(err)
	}

	for _, row := range rows {
		posts = append(posts, (*Post)(row))
	}
	return posts, nil
}
