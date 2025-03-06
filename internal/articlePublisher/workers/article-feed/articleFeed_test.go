package articlefeed

import (
	"testing"

	"github.com/ravilock/goduit/internal/app"
)

func TestArticleFeedWorker(t *testing.T) {
  articleGetterMock := newMockArticleGetter(t)
  profileGetterMock := newMockProfileGetter(t)
  followersGetterMock := newMockFollowersGetter(t)
  feedAppenderMock := newMockFeedAppender(t)
  articleWriteQueueConsumer := app.NewMockConsumer(t)
}
