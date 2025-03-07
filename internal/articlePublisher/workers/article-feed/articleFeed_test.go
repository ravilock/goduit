package articlefeed

import (
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ravilock/goduit/internal/app"
	"github.com/ravilock/goduit/internal/articlePublisher/models"
	followerCentralModels "github.com/ravilock/goduit/internal/followerCentral/models"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	mock "github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	articleWriteQueueConsumerMock := app.NewMockConsumer(t)
	worker := NewArticleFeedWorker(articleWriteQueueConsumerMock, articleGetterMock, profileGetterMock, followersGetterMock, feedAppenderMock, slog.New(logHandler))

	t.Run("Should receive new article message and append it to followers feed", func(t *testing.T) {
		// Arrange & Assert
		expectedAuthor := assembleUserModel()
		expectedAuthorID := *expectedAuthor.ID
		expectedArticle := assembleArticleModel(expectedAuthorID)
		expectedArticleID := expectedArticle.ID.Hex()
		expectedMessageBody := []byte(expectedArticleID)
		expectedFollowers := []*followerCentralModels.Follower{assembleFollowerModel(expectedArticleID), assembleFollowerModel(expectedArticleID)}
		followerIDs := make([]string, 0, len(expectedFollowers))
		for _, follower := range expectedFollowers {
			followerIDs = append(followerIDs, *follower.Follower)
		}
		messageMock := app.NewMockMessage(t)
		messageQueue := make(chan app.Message)
		defer close(messageQueue)
		articleWriteQueueConsumerMock.EXPECT().StartConsumer().Return().Once()
		articleWriteQueueConsumerMock.EXPECT().Consume().Return(messageQueue).Once()
		messageMock.EXPECT().Data().Return(expectedMessageBody).Once()
		articleGetterMock.EXPECT().GetArticleByID(mock.AnythingOfType("context.backgroundCtx"), expectedArticleID).Return(expectedArticle, nil).Once()
		profileGetterMock.EXPECT().GetUserByID(mock.AnythingOfType("context.backgroundCtx"), expectedAuthorID.Hex()).Return(expectedAuthor, nil).Once()
		followersGetterMock.EXPECT().GetFollowers(mock.AnythingOfType("context.backgroundCtx"), expectedAuthorID.Hex()).Return(expectedFollowers, nil).Once()
		feedAppenderMock.EXPECT().AppendArticleToUserFeeds(mock.AnythingOfType("context.backgroundCtx"), expectedArticle, expectedAuthor, followerIDs).Return(nil).Once()
		messageMock.EXPECT().Success().Return(nil).Once()

		// Act
		go worker.Consume()
		messageQueue <- messageMock
	})
}

func assembleArticleModel(authorID primitive.ObjectID) *models.Article {
	articleID := primitive.NewObjectID()
	authorIDHex := authorID.Hex()
	articleTitle := "Article Title"
	articleSlug := "article-title"
	articleDescription := "Article Description"
	articleBody := "Article Body"
	articleTagList := []string{"test"}
	now := time.Now().UTC().Truncate(time.Millisecond)
	return &models.Article{
		ID:             &articleID,
		Author:         &authorIDHex,
		Slug:           &articleSlug,
		Title:          &articleTitle,
		Description:    &articleDescription,
		Body:           &articleBody,
		TagList:        articleTagList,
		CreatedAt:      &now,
		UpdatedAt:      &now,
		FavoritesCount: new(int64),
	}
}

func assembleUserModel() *profileManagerModels.User {
	authorID := primitive.NewObjectID()
	username := "raylok"
	bio := "This is a bio"
	image := "https://cataas.com/cat"
	return &profileManagerModels.User{
		ID:       &authorID,
		Username: &username,
		Bio:      &bio,
		Image:    &image,
	}
}

func assembleFollowerModel(followed string) *followerCentralModels.Follower {
	ID := primitive.NewObjectID()
	followerID := uuid.NewString()
	return &followerCentralModels.Follower{
		ID:       &ID,
		Followed: &followed,
		Follower: &followerID,
	}
}
