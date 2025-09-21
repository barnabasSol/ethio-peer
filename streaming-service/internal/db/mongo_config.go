package db

import (
	"context"
	"ep-streaming-service/internal/models"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConfigProp(client *mongo.Client) {
	collection := client.Database(Name).Collection(models.SessionCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	field_index := mongo.IndexModel{
		Keys: bson.D{
			{Key: "owner_id", Value: 1},
			{Key: "owner_username", Value: 1},
			{Key: "session_name", Value: 1},
			{Key: "is_live", Value: 1},
		},
		Options: options.Index().SetUnique(false),
	}

	if name, err := collection.Indexes().CreateOne(
		ctx,
		field_index,
	); err != nil {
		log.Fatalf("Failed to create index: %v", err)
	} else {
		fmt.Println("Created session field index:", name)
	}

}
