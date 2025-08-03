package profiles

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
	"github.com/Archisman-Mridha/x-clone/backend/protobuf/generated"
)

type ProfilesAPI struct {
	generated.UnimplementedProfilesAPIServer

	profilesService *ProfilesService
}

func NewProfilesAPI(profilesService *ProfilesService) *ProfilesAPI {
	//nolint:exhaustruct
	return &ProfilesAPI{
		profilesService: profilesService,
	}
}

func (*ProfilesAPI) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *ProfilesAPI) SearchProfiles(ctx context.Context,
	request *generated.SearchProfilesRequest,
) (*generated.SearchProfilesResponse, error) {
	profilePreviews, err := p.profilesService.SearchProfiles(ctx, &SearchProfilesArgs{
		Query: request.GetQuery(),

		PageRequest: &utils.PageRequest{
			Offset: request.GetPageRequest().GetOffset(),
			Limit:  request.GetPageRequest().GetLimit(),
		},
	})
	if err != nil {
		return nil, err
	}

	response := &generated.SearchProfilesResponse{
		ProfilePreviews: toProtoGeneratedProfilePreviews(profilePreviews),
	}
	return response, nil
}

func (p *ProfilesAPI) GetProfilePreviews(ctx context.Context,
	request *generated.GetProfilePreviewsRequest,
) (*generated.GetProfilePreviewsResponse, error) {
	profilePreviews, err := p.profilesService.GetProfilePreviews(ctx, request.GetIds())
	if err != nil {
		return nil, err
	}

	response := &generated.GetProfilePreviewsResponse{
		ProfilePreviews: toProtoGeneratedProfilePreviews(profilePreviews),
	}
	return response, nil
}

// Converts []*ProfilePreview to []*generated.ProfilePreiew.
func toProtoGeneratedProfilePreviews(input []*ProfilePreview) []*generated.ProfilePreview {
	output := make([]*generated.ProfilePreview, len(input))
	for _, item := range input {
		output = append(output, &generated.ProfilePreview{
			Id:       item.ID,
			Name:     item.Name,
			Username: item.Username,
		})
	}
	return output
}
