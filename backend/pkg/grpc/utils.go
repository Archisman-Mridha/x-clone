package grpc

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

const GRPC_METADATA_USER_ID = "user_id"

var (
	ErrUserIDMetadataNotFound = utils.NewAPIError("user_id not found in gRPC request metadata")
	ErrParsingUserID          = utils.NewAPIError("failed parsing user_id in gRPC request metadata")
)

// Extracts and returns the user ID, from the gRPC request metadata.
func GetUserIDFromMetadata(ctx context.Context) (int32, utils.APIError) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, ErrUserIDMetadataNotFound
	}

	userIDs := metadata.Get(GRPC_METADATA_USER_ID)
	if len(userIDs) == 0 {
		return 0, ErrUserIDMetadataNotFound
	}

	userID := userIDs[0]

	parsedUserID, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		return 0, ErrParsingUserID
	}
	return int32(parsedUserID), nil
}
