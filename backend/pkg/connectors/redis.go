package connectors

import (
	"context"
	"log/slog"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/assert"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/logger"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type (
	RedisConnector struct {
		client redis.UniversalClient
	}

	NewRedisConnectorArgs struct {
		NodeAddresses []string

		Username,
		Password string
	}
)

func NewRedisConnector(ctx context.Context, args *NewRedisConnectorArgs) *RedisConnector {
	var client redis.UniversalClient
	switch {
	case len(args.NodeAddresses) > 1:
		//nolint:exhaustruct
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: args.NodeAddresses,

			Username: args.Username,
			Password: args.Password,
		})

	default:
		//nolint:exhaustruct
		client = redis.NewClient(&redis.Options{
			Addr: args.NodeAddresses[0],

			Username: args.Username,
			Password: args.Password,
		})
	}

	// Ping the cluster, verifying that a working connection has been established.
	err := client.Ping(ctx).Err()
	assert.AssertErrNil(ctx, err, "Failed connecting to Redis")

	slog.DebugContext(ctx, "Connected to Redis")

	// Instrument the Redis client
	{
		err = redisotel.InstrumentTracing(client)
		assert.AssertErrNil(ctx, err, "Failed trace instrumenting the Redis client")

		err = redisotel.InstrumentMetrics(client)
		assert.AssertErrNil(ctx, err, "Failed metric instrumenting the Redis client")
	}

	return &RedisConnector{client}
}

func (r *RedisConnector) GetClient() redis.UniversalClient {
	return r.client
}

func (r *RedisConnector) Healthcheck() error {
	if err := r.client.Ping(context.Background()).Err(); err != nil {
		return utils.WrapErrorWithPrefix("Failed pinging Redis cluster", err)
	}
	return nil
}

func (r *RedisConnector) Shutdown() {
	if err := r.client.Close(); err != nil {
		slog.Error("Failed closing Redis connection", logger.Error(err))
		return
	}
	slog.Debug("Shut down Redis client")
}
