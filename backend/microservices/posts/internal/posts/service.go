package posts

import (
	"context"

	goValidator "github.com/go-playground/validator/v10"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

type PostsService struct {
	validator *goValidator.Validate

	postsRepository PostsRepository
}

func NewPostsService(
	validator *goValidator.Validate,
	postsRespository PostsRepository,
) *PostsService {
	return &PostsService{
		validator,
		postsRespository,
	}
}

type CreatePostArgs struct {
	OwnerID     int32
	Description string `validate:"description"`
}

func (p *PostsService) CreatePost(ctx context.Context, args *CreatePostArgs) (int32, error) {
	// Validate input.
	err := p.validator.StructCtx(ctx, args)
	if err != nil {
		return 0, err
	}

	return p.postsRepository.Create(ctx, args)
}

type GetPostsOfUserArgs struct {
	OwnerID     int32
	PageRequest *utils.PageRequest
}

func (p *PostsService) GetUserPosts(
	ctx context.Context,
	args *GetPostsOfUserArgs,
) ([]*Post, error) {
	return p.postsRepository.GetPostsOfUser(ctx, args)
}

func (p *PostsService) GetPosts(ctx context.Context, ids []int32) ([]*Post, error) {
	return p.postsRepository.GetPosts(ctx, ids)
}
