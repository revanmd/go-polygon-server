package database

import (
	"context"
	"log"
	"time"

	"polygon-server/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(cfg *config.Config) *mongo.Database {
	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := client.Database(cfg.DBName)
	log.Println("Connected to MongoDB!")
	return db
}
