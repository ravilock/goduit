package requests

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/stretchr/testify/require"
)

func TestFeedArticles(t *testing.T) {
	t.Run("Valid request should return errors", func(t *testing.T) {
		request := generateFeedArticlesRequest()
		err := request.Validate()
		require.NoError(t, err)
	})
	t.Run("Limit should have min value 1", func(t *testing.T) {
		request := generateFeedArticlesRequest()
		request.Pagination.Limit = 0
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Limit", "min", "1").Error())
	})
	t.Run("Limit should have max value 30", func(t *testing.T) {
		request := generateFeedArticlesRequest()
		request.Pagination.Limit = 31
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Limit", "max", "30").Error())
	})
	t.Run("Offset should have min value 1", func(t *testing.T) {
		request := generateFeedArticlesRequest()
		request.Pagination.Offset = -1
		err := request.Validate()
		require.ErrorContains(t, err, api.InvalidFieldLimit("Offset", "min", "0").Error())
	})
}

func generateFeedArticlesRequest() *FeedArticlesRequest {
	return &FeedArticlesRequest{
		FeedArticlesPagination{
			Limit:  20,
			Offset: 20,
		},
	}
}
