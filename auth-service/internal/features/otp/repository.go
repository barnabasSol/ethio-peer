package otp

import (
	"context"
	"ep-auth-service/internal/db"
	"ep-auth-service/internal/features/shared"
	"ep-auth-service/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	GetUserById(ctx context.Context, user_id string) (*models.User, error)
}
type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUserById(
	ctx context.Context,
	user_id string,
) (*models.User, error) {
	user_collection := r.db.Database(db.Name).Collection(models.UserCollection)
	filter := bson.D{{Key: "_id", Value: user_id}}
	var user models.User
	err := user_collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, shared.ErrUserNotFound
	}
	return &user, nil
}
