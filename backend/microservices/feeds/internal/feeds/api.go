package feeds

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/grpc"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
	"github.com/Archisman-Mridha/x-clone/backend/protobuf/generated"
)

type FeedsAPI struct {
	generated.UnimplementedFeedsAPIServer

	feedsService *FeedsService
}

func NewFeedsAPI(feedsService *FeedsService) *FeedsAPI {
	//nolint:exhaustruct
	return &FeedsAPI{
		feedsService: feedsService,
	}
}

func (*FeedsAPI) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (f *FeedsAPI) GetFeed(ctx context.Context,
	request *generated.GetFeedRequest,
) (*generated.GetFeedResponse, error) {
	userID, err := grpc.GetUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	postIDs, err := f.feedsService.GetFeed(ctx, &GetFeedArgs{
		UserID: userID,

		PageRequest: &utils.PageRequest{
			Offset: request.GetPageRequest().GetOffset(),
			Limit:  request.GetPageRequest().GetLimit(),
		},
	})
	if err != nil {
		return nil, err
	}

	response := &generated.GetFeedResponse{
		PostIds: postIDs,
	}
	return response, nil
}
