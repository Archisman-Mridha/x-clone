package followships

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/grpc"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
	"github.com/Archisman-Mridha/x-clone/backend/protobuf/generated"
)

type FollowshipsAPI struct {
	generated.UnimplementedFollowshipsAPIServer

	followshipsService *FollowshipsService
}

func NewFollowshipsAPI(followshipsService *FollowshipsService) *FollowshipsAPI {
	//nolint:exhaustruct
	return &FollowshipsAPI{
		followshipsService: followshipsService,
	}
}

func (*FollowshipsAPI) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (f *FollowshipsAPI) Follow(ctx context.Context,
	request *generated.FollowRequest,
) (*emptypb.Empty, error) {
	followerID, err := grpc.GetUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	err = f.followshipsService.CreateFollowship(ctx, &FollowshipOperationArgs{
		FollowerID: followerID,
		FolloweeID: request.GetFolloweeId(),
	})
	return &emptypb.Empty{}, err
}

func (f *FollowshipsAPI) Unfollow(ctx context.Context,
	request *generated.UnfollowRequest,
) (*emptypb.Empty, error) {
	followerID, err := grpc.GetUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	err = f.followshipsService.DeleteFollowship(ctx, &FollowshipOperationArgs{
		FollowerID: followerID,
		FolloweeID: request.GetFolloweeId(),
	})
	return &emptypb.Empty{}, err
}

func (f *FollowshipsAPI) Follows(ctx context.Context,
	request *generated.FollowsRequest,
) (*generated.FollowsResponse, error) {
	followerID, err := grpc.GetUserIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	followshipExists, err := f.followshipsService.FollowshipExists(ctx, &FollowshipOperationArgs{
		FollowerID: followerID,
		FolloweeID: request.GetFolloweeId(),
	})
	if err != nil {
		return nil, err
	}

	response := &generated.FollowsResponse{
		Follows: followshipExists,
	}
	return response, err
}

func (f *FollowshipsAPI) GetFollowers(ctx context.Context,
	request *generated.GetFollowersRequest,
) (*generated.GetFollowersResponse, error) {
	followerIDs, err := f.followshipsService.GetFollowers(ctx, &GetFollowersArgs{
		FolloweeID: request.GetUserId(),

		PageRequest: &utils.PageRequest{
			Offset: request.GetPageRequest().GetOffset(),
			Limit:  request.GetPageRequest().GetLimit(),
		},
	})
	if err != nil {
		return nil, err
	}

	response := &generated.GetFollowersResponse{
		FollowerIds: followerIDs,
	}
	return response, nil
}

func (f *FollowshipsAPI) GetFollowees(ctx context.Context,
	request *generated.GetFolloweesRequest,
) (*generated.GetFolloweesResponse, error) {
	followeeIDs, err := f.followshipsService.GetFollowees(ctx, &GetFolloweesArgs{
		FollowerID: request.GetUserId(),

		PageRequest: &utils.PageRequest{
			Offset: request.GetPageRequest().GetOffset(),
			Limit:  request.GetPageRequest().GetLimit(),
		},
	})
	if err != nil {
		return nil, err
	}

	response := &generated.GetFolloweesResponse{
		FolloweeIds: followeeIDs,
	}
	return response, nil
}

func (f *FollowshipsAPI) GetFollowerAndFolloweeCounts(ctx context.Context,
	request *generated.GetFollowerAndFolloweeCountsRequest,
) (*generated.GetFollowerAndFolloweeCountsResponse, error) {
	followerAndFolloweeCounts, err := f.followshipsService.GetFollowerAndFolloweeCounts(ctx,
		request.GetUserId(),
	)
	if err != nil {
		return nil, err
	}

	response := &generated.GetFollowerAndFolloweeCountsResponse{
		FollowerCount: followerAndFolloweeCounts.FollowerCount,
		FolloweeCount: followerAndFolloweeCounts.FolloweeCount,
	}
	return response, nil
}
