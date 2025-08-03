package profiles

import (
	"context"

	goValidator "github.com/go-playground/validator/v10"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

type ProfilesService struct {
	validator *goValidator.Validate

	profilesRepository   ProfilesRepository
	profilesSearchEngine ProfilesSearchEngine
}

func NewProfilesService(
	validator *goValidator.Validate,
	profilesRespository ProfilesRepository,
	profilesSearchEngine ProfilesSearchEngine,
) *ProfilesService {
	return &ProfilesService{
		validator,
		profilesRespository,
		profilesSearchEngine,
	}
}

type CreateProfileArgs struct {
	ID int32
	Name,
	Username string
}

func (p *ProfilesService) CreateProfile(ctx context.Context, args *CreateProfileArgs) error {
	return p.profilesRepository.Create(ctx, args)
}

func (p *ProfilesService) GetProfilePreviews(ctx context.Context,
	ids []int32,
) ([]*ProfilePreview, error) {
	return p.profilesRepository.GetPreviews(ctx, ids)
}

func (p *ProfilesService) IndexProfile(ctx context.Context, profilePreview *ProfilePreview) error {
	return p.profilesSearchEngine.IndexProfile(ctx, profilePreview)
}

type SearchProfilesArgs struct {
	Query       string `validate:"notblank"`
	PageRequest *utils.PageRequest
}

func (p *ProfilesService) SearchProfiles(ctx context.Context,
	args *SearchProfilesArgs,
) ([]*ProfilePreview, error) {
	// Validate input.
	err := p.validator.StructCtx(ctx, args)
	if err != nil {
		return nil, err
	}

	return p.profilesSearchEngine.SearchProfiles(ctx, (*SearchProfilesArgs)(args))
}
