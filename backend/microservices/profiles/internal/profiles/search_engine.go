package profiles

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/aquasecurity/esquery"

	"github.com/Archisman-Mridha/x-clone/backend/microservices/profiles/internal/constants"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/assert"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/connectors"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"
)

type ProfilesSearchEngine interface {
	// This operation needs to be idempotent, since it's invoked by a DB event processor.
	IndexProfile(ctx context.Context, profilePreview *ProfilePreview) error

	SearchProfiles(ctx context.Context, args *SearchProfilesArgs) ([]*ProfilePreview, error)
}

type (
	ProfilesElasticsearchSearchEngine struct {
		*connectors.ElasticsearchConnector
	}

	ProfileMetadata struct {
		Name     string `json:"name,omitempty"`
		Username string `json:"username,omitempty"`
	}
)

func NewProfilesElasticsearchSearchEngine(
	elasticsearchConnector *connectors.ElasticsearchConnector,
) *ProfilesElasticsearchSearchEngine {
	return &ProfilesElasticsearchSearchEngine{elasticsearchConnector}
}

func (s *ProfilesElasticsearchSearchEngine) IndexProfile(ctx context.Context,
	profilePreview *ProfilePreview,
) error {
	profileMetadata := ProfileMetadata{
		Name:     profilePreview.Name,
		Username: profilePreview.Username,
	}
	jsonEncodedProfileMetadata, err := json.Marshal(profileMetadata)
	if err != nil {
		return utils.WrapError(err)
	}

	elasticsearchIndexClient := s.GetClient().Index

	response, err := elasticsearchIndexClient(
		constants.SEARCH_ENGINE_INDEX_PROFILES,
		bytes.NewReader(jsonEncodedProfileMetadata),
		elasticsearchIndexClient.WithDocumentID(string(profilePreview.ID)),
		elasticsearchIndexClient.WithContext(ctx),
	)
	if err != nil {
		return utils.WrapError(err)
	} else if response.IsError() {
		return utils.WrapError(errors.New("failed indexing profile, received error status-code"))
	}
	//nolint:errcheck
	defer response.Body.Close()

	return nil
}

func (s *ProfilesElasticsearchSearchEngine) SearchProfiles(ctx context.Context,
	args *SearchProfilesArgs,
) ([]*ProfilePreview, error) {
	profilePreviews := []*ProfilePreview{}

	searchQuery, err := esquery.Search().
		Query(
			esquery.MultiMatch(args.Query).
				Type(esquery.MatchTypePhrasePrefix).
				Fields("name", "username"),
		).
		Sort("_id", esquery.OrderAsc).
		SearchAfter(args.PageRequest.Offset).
		Size(uint64(args.PageRequest.Limit)).
		MarshalJSON()
	assert.AssertErrNil(ctx, err, "Failed JSON marshalling Elasticsearch search query")

	elasticsearchSearchClient := s.GetClient().Search

	response, err := elasticsearchSearchClient(
		elasticsearchSearchClient.WithBody(bytes.NewReader(searchQuery)),
		elasticsearchSearchClient.WithContext(ctx),
		elasticsearchSearchClient.WithIndex(constants.SEARCH_ENGINE_INDEX_PROFILES),
	)
	if (err != nil) || (response.IsError()) {
		return profilePreviews, utils.WrapErrorWithPrefix(
			"Failed running Elasticsearch search query",
			err,
		)
	}
	//nolint:errcheck
	defer response.Body.Close()

	parsedResponseBody, err := utils.ParseElasticsearchSearchQueryResponseBody[ProfileMetadata](ctx,
		response.Body,
	)
	if err != nil {
		return profilePreviews, utils.WrapError(err)
	}

	for _, hit := range parsedResponseBody.Hits.Hits {
		id, err := strconv.ParseInt(hit.ID, 10, 32)
		if err != nil {
			return profilePreviews, utils.WrapErrorWithPrefix(
				"Failed parsing profile id to int32",
				err,
			)
		}

		profilePreviews = append(profilePreviews, &ProfilePreview{
			ID:       int32(id),
			Name:     hit.Source.Name,
			Username: hit.Source.Username,
		})
	}

	return profilePreviews, nil
}
