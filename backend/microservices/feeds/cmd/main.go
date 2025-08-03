package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"

	"github.com/Archisman-Mridha/x-clone/backend/microservices/feeds/internal/config"
	"github.com/Archisman-Mridha/x-clone/backend/microservices/feeds/internal/constants"
	"github.com/Archisman-Mridha/x-clone/backend/microservices/feeds/internal/feeds"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/connectors"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/grpc"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/healthcheck"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/logger"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
	"github.com/Archisman-Mridha/x-clone/backend/protobuf/generated"
)

var configFilePath string

func main() {
	// Get CLI flag values.
	{
		flagSet := flag.NewFlagSet("", flag.ExitOnError)

		flagSet.StringVar(&configFilePath,
			constants.FLAG_CONFIG_FILE, constants.FLAG_CONFIG_FILE_DEFAULT,
			"Config file path",
		)

		flagSet.VisitAll(utils.CreateGetFlagOrEnvValueFn(""))

		cmdArgs := os.Args[1:]
		if err := flagSet.Parse(cmdArgs); err != nil {
			slog.Error("Failed parsing command line flags", logger.Error(err))
			os.Exit(1)
		}
	}

	// When the program receives any interruption / SIGKILL / SIGTERM signal, the cancel function is
	// automatically invoked. The cancel function is responsible for freeing all the resources
	// associated with the context.
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)

	// Construct validator with custom validators.
	validator := utils.NewValidator(ctx)

	// Get config.
	config := utils.MustParseConfigFile[config.Config](ctx, configFilePath, validator)

	if err := run(ctx, config); err != nil {
		slog.ErrorContext(ctx, err.Error())

		cancel()

		// Give some time for remaining resources (if any) to be cleaned up.
		time.Sleep(time.Second)

		os.Exit(1)
	}
}

func run(ctx context.Context, config *config.Config) error {
	waitGroup, ctx := errgroup.WithContext(ctx)

	// Construct connectors.

	redisConnector := connectors.NewRedisConnector(ctx,
		(*connectors.NewRedisConnectorArgs)(&config.Redis),
	)
	defer redisConnector.Shutdown()

	// Construct services.

	feedsService := feeds.NewFeedsService(
		feeds.NewFeedsRedisRepository(redisConnector),
	)

	// Run gRPC server.

	gRPCServer := grpc.NewGRPCServer(ctx, grpc.NewGRPCServerArgs{
		DevModeEnabled: config.DevMode,

		Healthcheckables: []healthcheck.Healthcheckable{
			redisConnector,
		},

		ToGRPCErrorStatusCodeFn: getGRPCErrorStatusCode,
	})

	feedsAPI := feeds.NewFeedsAPI(feedsService)
	generated.RegisterFeedsAPIServer(gRPCServer, feedsAPI)

	waitGroup.Go(func() error {
		return gRPCServer.Run(ctx, config.ServerPort)
	})

	/*
		The returned channel gets closed when either of this happens :

			(1) A program termination signal is received, because of which the parent context's done
				  channel gets closed.

			(2) Any of the go-routines registered under the wait-group, finishes running.
	*/
	<-ctx.Done()
	slog.DebugContext(ctx, "Gracefully shutting down program")

	// Gracefully shutdown the gRPC server, ensuring that it finishes ongoing processing of requests.
	gRPCServer.GracefulShutdown()

	return waitGroup.Wait()
}

// Returns suitable gRPC error status code, based on the given error.
func getGRPCErrorStatusCode(err error) codes.Code {
	_, ok := err.(utils.APIError)
	if !ok {
		return codes.Internal
	}

	return codes.Unknown
}
