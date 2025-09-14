package db

import (
	"context"
	"ep-peer-service/internal/models"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConfigProp(client *mongo.Client) {
	collection := client.Database(Name).Collection(models.PeerCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	onlineStatusScoreIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "online_status", Value: 1},
			{Key: "overall_score", Value: 1},
		},
		Options: options.Index().SetUnique(false),
	}

	if name, err := collection.Indexes().CreateOne(
		ctx,
		onlineStatusScoreIndex,
	); err != nil {
		log.Fatalf("Failed to create index: %v", err)
	} else {
		fmt.Println("Created online status index:", name)
	}

}
