package articlefeed

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/ravilock/goduit/internal/app"
)

type loggerSpy struct {
	LastMessage   string
	NumberOfCalls int
}

func (l *loggerSpy) Write(p []byte) (int, error) {
	l.NumberOfCalls++
	l.LastMessage = string(p)
	return len(p), nil
}

func (l *loggerSpy) Clean() {
	l.NumberOfCalls = 0
	l.LastMessage = ""
}

func TestArticleFeedWorker(t *testing.T) {
	logSpy := new(loggerSpy)
	logHandler := slog.NewTextHandler(logSpy, nil)
	articleGetterMock := newMockArticleGetter(t)
	profileGetterMock := newMockProfileGetter(t)
	followersGetterMock := newMockFollowersGetter(t)
	feedAppenderMock := newMockFeedAppender(t)
	articleWriteQueueConsumer := app.NewMockConsumer(t)
	worker := NewArticleFeedWorker(articleWriteQueueConsumer, articleGetterMock, profileGetterMock, followersGetterMock, feedAppenderMock, slog.New(logHandler))

	t.Run("Should receive new article message and append it to followers feed", func(t *testing.T) {
		fmt.Println(worker)
	})
}
