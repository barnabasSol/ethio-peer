package otp

import (
	"context"
	"ep-auth-service/internal/db"
	"ep-auth-service/internal/features/common"
	"ep-auth-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	GetUserById(ctx context.Context, user_id string) (*models.User, error)
	UpdateUser(
		ctx context.Context,
		user_id bson.ObjectID,
		email_verified, is_active bool,
	) (bool, error)
}
type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) UpdateUser(
	ctx context.Context,
	user_id bson.ObjectID,
	email_verified, is_active bool,
) (bool, error) {
	user_collection := r.db.Database(db.Name).Collection(models.UserCollection)

	filter := bson.D{{Key: "_id", Value: user_id}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "institute_email_verified", Value: email_verified},
			{Key: "is_active", Value: is_active},
			{Key: "updated_at", Value: time.Now().UTC()},
		}},
	}

	updateResult, err := user_collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	if updateResult.MatchedCount == 0 {
		return false, common.ErrUserNotFound
	}

	return true, nil
}

func (r *repository) GetUserById(
	ctx context.Context,
	user_id string,
) (*models.User, error) {
	user_collection := r.db.Database(db.Name).Collection(models.UserCollection)
	user_obj_id, err := bson.ObjectIDFromHex(user_id)
	if err != nil {
		return nil, common.ErrUserNotFound
	}
	filter := bson.D{{Key: "_id", Value: user_obj_id}}
	var user models.User
	err = user_collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, common.ErrUserNotFound
	}
	return &user, nil
}
