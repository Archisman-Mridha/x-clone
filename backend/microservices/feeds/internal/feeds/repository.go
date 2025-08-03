package feeds

import (
	"context"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/connectors"
)

type (
	FeedsRepository interface {
		/*
			Returns the specified elements of the list stored at key.

			The offsets start and stop are zero-based indexes.
			These offsets can also be negative numbers indicating offsets starting at the end of the
			list.

			Out of range indexes will not produce an error :

			  (1) If start is larger than the end of the list, an empty list is returned.

			  (2) If stop is larger than the actual end of the list, Redis will treat it like the last
			      element of the list.
		*/
		LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	}
)

type FeedsRedisRepository struct {
	redisConnector *connectors.RedisConnector
}

func NewFeedsRedisRepository(redisConnector *connectors.RedisConnector) FeedsRepository {
	return &FeedsRedisRepository{
		redisConnector,
	}
}

func (f *FeedsRedisRepository) LRange(ctx context.Context,
	key string,
	start, stop int64,
) ([]string, error) {
	return f.redisConnector.GetClient().LRange(ctx, key, start, stop).Result()
}
