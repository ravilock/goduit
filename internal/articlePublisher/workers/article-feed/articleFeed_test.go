package articlefeed

import (
	"testing"
)

func TestArticleFeedWorker(t *testing.T) {
  articleGetterMock := newMockArticleGetter(t)
  profileGetterMock := newMockProfileGetter(t)
  followersGetterMock := newMockFollowersGetter(t)
  feedAppenderMock := newMockFeedAppender(t)
}
