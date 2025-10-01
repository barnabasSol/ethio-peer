package user

import (
	"context"
	"ep-auth-service/internal/db"
	"ep-auth-service/internal/models"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	GetUser(
		ctx context.Context,
		user_id bson.ObjectID,
	) (*models.User, error)
}

type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUser(
	ctx context.Context,
	user_id bson.ObjectID,
) (*models.User, error) {
	user_collection := r.db.Database(db.Name).Collection(models.UserCollection)
	filter := bson.M{"_id": user_id}
	var user models.User

	err := user_collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
