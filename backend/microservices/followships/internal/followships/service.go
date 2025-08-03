package followships

import (
	"context"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

type FollowshipsService struct {
	followshipsRepository FollowshipsRepository
}

func NewFollowshipsService(followshipsRespository FollowshipsRepository) *FollowshipsService {
	return &FollowshipsService{
		followshipsRespository,
	}
}

type FollowshipOperationArgs struct {
	FollowerID,
	FolloweeID int32
}

func (f *FollowshipsService) CreateFollowship(ctx context.Context,
	args *FollowshipOperationArgs,
) error {
	return f.followshipsRepository.CreateFollowship(ctx, args)
}

func (f *FollowshipsService) DeleteFollowship(ctx context.Context,
	args *FollowshipOperationArgs,
) error {
	return f.followshipsRepository.DeleteFollowship(ctx, args)
}

func (f *FollowshipsService) FollowshipExists(ctx context.Context,
	args *FollowshipOperationArgs,
) (bool, error) {
	return f.followshipsRepository.FollowshipExists(ctx, args)
}

type GetFollowersArgs struct {
	FolloweeID int32

	PageRequest *utils.PageRequest
}

func (f *FollowshipsService) GetFollowers(ctx context.Context,
	args *GetFollowersArgs,
) ([]int32, error) {
	return f.followshipsRepository.GetFollowers(ctx, args)
}

type GetFolloweesArgs struct {
	FollowerID int32

	PageRequest *utils.PageRequest
}

func (f *FollowshipsService) GetFollowees(ctx context.Context,
	args *GetFolloweesArgs,
) ([]int32, error) {
	return f.followshipsRepository.GetFollowees(ctx, args)
}

type GetFollowerAndFolloweeCountsOutput struct {
	FollowerCount,
	FolloweeCount int64
}

func (f *FollowshipsService) GetFollowerAndFolloweeCounts(ctx context.Context,
	userID int32,
) (*GetFollowerAndFolloweeCountsOutput, error) {
	return f.followshipsRepository.GetFollowerAndFolloweeCounts(ctx, userID)
}
