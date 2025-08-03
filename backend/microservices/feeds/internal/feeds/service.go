package feeds

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

type FeedsService struct {
	feedsRepository FeedsRepository
}

func NewFeedsService(feedsRepository FeedsRepository) *FeedsService {
	return &FeedsService{
		feedsRepository,
	}
}

type GetFeedArgs struct {
	UserID int32

	PageRequest *utils.PageRequest
}

func (f *FeedsService) GetFeed(ctx context.Context, args *GetFeedArgs) ([]int32, error) {
	output, err := f.feedsRepository.LRange(ctx,
		string(args.UserID),
		int64(args.PageRequest.Offset),
		int64(args.PageRequest.Offset+args.PageRequest.Limit-1),
	)
	if err != nil {
		return []int32{}, utils.WrapError(err)
	}

	postIDs := []int32{}
	for _, item := range output {
		postID, err := strconv.Atoi(item)
		if err != nil {
			slog.WarnContext(ctx, "Failed parsing post ID", slog.String("value", item))
			continue
		}

		postIDs = append(postIDs, int32(postID))
	}
	return postIDs, nil
}
