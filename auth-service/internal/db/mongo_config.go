package db

import (
	"context"
	"ep-auth-service/internal/models"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConfigProp(client *mongo.Client) {
	collection := client.Database(Name).Collection(models.UserCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if name, err := collection.Indexes().CreateOne(ctx, usernameIndex); err != nil {
		log.Fatalf("Failed to create username index: %v", err)
	} else {
		fmt.Println("Created username index:", name)
	}

	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetSparse(true), // sparse allows missing email
	}
	if name, err := collection.Indexes().CreateOne(ctx, emailIndex); err != nil {
		log.Fatalf("Failed to create email index: %v", err)
	} else {
		fmt.Println("Created email index:", name)
	}

	instituteEmailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "institute_email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if name, err := collection.Indexes().CreateOne(ctx, instituteEmailIndex); err != nil {
		log.Fatalf("Failed to create institute_email index: %v", err)
	} else {
		fmt.Println("Created institute_email index:", name)
	}
}
