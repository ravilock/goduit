package handlers_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ravilock/goduit/internal/api"
	"github.com/ravilock/goduit/internal/config"
	"github.com/ravilock/goduit/internal/identity"
	profileManagerModels "github.com/ravilock/goduit/internal/profileManager/models"
	profileManager "github.com/ravilock/goduit/internal/profileManager/services"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	viper.SetDefault("server.url", "http://localhost:9090")
	privateKeyFile, err := os.Open(viper.GetString("private.key.location"))
	if err != nil {
		log.Fatal("Failed to open private key file", err)
	}

	if err := config.LoadPrivateKey(privateKeyFile); err != nil {
		log.Fatal("Failed to load private key file content", err)
	}

	if err := privateKeyFile.Close(); err != nil {
		log.Fatal("Failed to close private key file", err)
	}

	publicKeyFile, err := os.Open(viper.GetString("public.key.location"))
	if err != nil {
		log.Fatal("Failed to open public key file", err)
	}

	if err := config.LoadPublicKey(publicKeyFile); err != nil {
		log.Fatal("Failed to load public key file content", err)
	}

	if err := publicKeyFile.Close(); err != nil {
		log.Fatal("Failed to close publicKeyFile key file", err)
	}

	server, err := api.NewServer()
	if err != nil {
		log.Fatalln("Failed to start server", err)
	}
	server.Start()
}

func clearDatabase(client *mongo.Client) {
	conduitDb := client.Database("conduit")
	collections, err := conduitDb.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Fatal("Could not list collections", err)
	}
	for _, coll := range collections {
		_, err := conduitDb.Collection(coll).DeleteMany(context.Background(), bson.D{})
		if err != nil {
			log.Fatal("Could not clear database", err)
		}
	}
}

func registerUser(username, email, password string, manager *profileManager.ProfileManager) (*identity.Identity, error) {
	if username == "" {
		username = "default-username"
	}
	if email == "" {
		email = "default.email@test.test"
	}
	if password == "" {
		password = "default-password"
	}
	createdAt := time.Now().Truncate(time.Millisecond)
	token, err := manager.Register(context.Background(), &profileManagerModels.User{
		Username:    &username,
		Email:       &email,
		CreatedAt:   &createdAt,
		LastSession: &createdAt,
	}, password)
	if err != nil {
		return nil, err
	}
	return identity.FromToken(token)
}
