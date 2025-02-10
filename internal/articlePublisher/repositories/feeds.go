package repositories

import (
	"context"
	"fmt"

	"github.com/ravilock/goduit/internal/articlePublisher/models"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FeedRepository struct {
	DBClient *mongo.Client
}

func NewFeedRepository(client *mongo.Client) *FeedRepository {
	return &FeedRepository{client}
}

func (r *FeedRepository) AppendArticleToUserFeeds(ctx context.Context, article *models.Article, author *profileManagerModels.User, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	feeds, err := r.listFeedsByID(ctx, userIDs)
	if err != nil {
		return err
	}

	feedFragment := assembleFeedFragmentFromArticle(article)

	maxFeedArticles := viper.GetInt("feed.max.articles")
	for i := range feeds {
		feeds[i].Articles = append(feeds[i].Articles, feedFragment)
		if len(feeds[i].Articles) > maxFeedArticles {
			feeds[i].Articles = feeds[i].Articles[1 : maxFeedArticles+1]
		}
	}

	usersWithoutFeed, err := getUsersWithoutFeed(userIDs, feeds)
	if err != nil {
		return err
	}

	for _, userID := range usersWithoutFeed {
		feed := models.Feed{
			UserID:   &userID,
			Articles: []models.FeedFragment{feedFragment},
		}
		feeds = append(feeds, feed)
	}

	collection := r.DBClient.Database("conduit").Collection("feeds")
	for _, feed := range feeds {
		filter := bson.M{"_id": feed.UserID}
		update := bson.M{"$set": feed}
		opt := options.Update().SetUpsert(true)

		result, err := collection.UpdateOne(ctx, filter, update, opt)
		if err != nil {
			return err
		}
		if result.ModifiedCount+result.UpsertedCount != int64(len(userIDs)) {
			return fmt.Errorf("mismatched result count. modified: %d, upserted: %d, expected: %d", result.ModifiedCount, result.UpsertedCount, len(userIDs))
		}
	}
	return nil
}

func (r *FeedRepository) listFeedsByID(ctx context.Context, IDs []string) ([]models.Feed, error) {
	objectIDs := make([]primitive.ObjectID, 0, len(IDs))
	for _, ID := range IDs {
		objID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objID)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	results := make([]models.Feed, len(objectIDs))
	collection := r.DBClient.Database("conduit").Collection("feeds")

	cursor, err := collection.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return results, err
	}

	return results, nil
}

func getUsersWithoutFeed(IDs []string, feeds []models.Feed) ([]primitive.ObjectID, error) {
	usersWithoutFeed := []primitive.ObjectID{}

	feedIdMap := make(map[string]struct{})
	for _, feed := range feeds {
		feedIdMap[feed.UserID.Hex()] = struct{}{}
	}

	for _, ID := range IDs {
		objID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			return nil, err
		}
		if _, ok := feedIdMap[ID]; !ok {
			usersWithoutFeed = append(usersWithoutFeed, objID)
		}
	}

	return usersWithoutFeed, nil
}

func assembleFeedFragmentFromArticle(article *models.Article) models.FeedFragment {
	articleID := article.ID.Hex()
	feedFragment := models.FeedFragment{
		ArticleID: &articleID,
		Author:    article.Author,
	}
	return feedFragment
}
