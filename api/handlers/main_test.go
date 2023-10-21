package handlers

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ravilock/goduit/api/validators"
	encryptionkeys "github.com/ravilock/goduit/internal/config/encryptionKeys"
	"github.com/ravilock/goduit/internal/config/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	if err := os.Chdir("../.."); err != nil {
		log.Fatalln("Error chanigng directory", err)
	}

	if err := godotenv.Load(".env.test"); err != nil {
		log.Fatalln("No .env file found", err)
	}

	if err := encryptionkeys.LoadKeys(); err != nil {
		log.Fatalln("Failed to read encrpytion keys", err)
	}

	databaseURI := os.Getenv("DB_URL")
	if databaseURI == "" {
		log.Fatalln("You must sey your 'DATABASE_URI' environmental variable.")
	}
	// Connect Mongo DB
	if err := mongo.ConnectDatabase(databaseURI); err != nil {
		log.Fatalln("Error connecting to database", err)
	}

	// Start Validator
	validators.InitValidator()
}

func clearDatabase() {
	conduitDb := mongo.DatabaseClient.Database("conduit")
	collections, err := conduitDb.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Fatal("Could not list collections", err)
	}
	for _, coll := range collections {
		conduitDb.Collection(coll).DeleteMany(context.Background(), bson.D{})
	}
}
