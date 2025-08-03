package posts

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/grpc"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
	"github.com/Archisman-Mridha/x-clone/backend/protobuf/generated"
)

type PostsAPI struct {
	generated.UnimplementedPostsAPIServer

	postsService *PostsService
}

func NewPostsAPI(postsService *PostsService) *PostsAPI {
	//nolint:exhaustruct
	return &PostsAPI{
		postsService: postsService,
	}
}

func (*PostsAPI) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *PostsAPI) CreatePost(ctx context.Context,
	request *generated.CreatePostRequest,
) (*generated.CreatePostResponse, error) {
	ownerID, err := grpc.GetUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	postID, err := p.postsService.CreatePost(ctx, &CreatePostArgs{
		OwnerID:     ownerID,
		Description: request.GetDescription(),
	})
	if err != nil {
		return nil, err
	}

	response := &generated.CreatePostResponse{
		PostId: postID,
	}
	return response, nil
}

func (p *PostsAPI) GetPostsOfUser(ctx context.Context,
	request *generated.GetPostsOfUserRequest,
) (*generated.GetPostsResponse, error) {
	posts, err := p.postsService.GetUserPosts(ctx, &GetPostsOfUserArgs{
		OwnerID: request.GetOwnerId(),

		PageRequest: &utils.PageRequest{
			Offset: request.GetPageRequest().GetOffset(),
			Limit:  request.GetPageRequest().GetLimit(),
		},
	})
	if err != nil {
		return nil, err
	}

	response := &generated.GetPostsResponse{
		Posts: toProtoGeneratedPosts(posts),
	}
	return response, nil
}

func (p *PostsAPI) GetPosts(ctx context.Context,
	request *generated.GetPostsRequest,
) (*generated.GetPostsResponse, error) {
	posts, err := p.postsService.GetPosts(ctx, request.GetPostIds())
	if err != nil {
		return nil, err
	}

	response := &generated.GetPostsResponse{
		Posts: toProtoGeneratedPosts(posts),
	}
	return response, nil
}

// Converts []*Post to []*generated.Post.
func toProtoGeneratedPosts(input []*Post) []*generated.Post {
	output := make([]*generated.Post, len(input))
	for _, item := range input {
		output = append(output, &generated.Post{
			Id:          item.ID,
			OwnerId:     item.OwnerID,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.String(),
		})
	}
	return output
}
