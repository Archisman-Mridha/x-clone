package utils

import (
	"context"
	"encoding/json"
	"io"
)

type (
	ResponseBody[T any] struct {
		Hits Hits[T] `json:"hits"`
	}

	Hits[T any] struct {
		Hits []Hit[T] `json:"hits"`
	}

	Hit[T any] struct {
		ID     string `json:"_id"`
		Source T      `json:"_source"`
	}
)

func ParseElasticsearchSearchQueryResponseBody[T any](ctx context.Context,
	responseBodyReader io.ReadCloser,
) (*ResponseBody[T], error) {
	responseBody := []byte{}
	_, err := responseBodyReader.Read(responseBody)
	if err != nil {
		return nil, WrapErrorWithPrefix(
			"Failed reading Elasticsearch search query response body", err,
		)
	}

	parsedResponseBody := new(ResponseBody[T])
	err = json.Unmarshal(responseBody, parsedResponseBody)
	if err != nil {
		return nil, WrapErrorWithPrefix(
			"Failed parsing Elasticsearch search query response body",
			err,
		)
	}
	return parsedResponseBody, nil
}
