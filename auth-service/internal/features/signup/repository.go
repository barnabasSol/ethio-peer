package signup

import (
	"context"
	"ep-auth-service/internal/db"
	"ep-auth-service/internal/models"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	InsertRefreshToken(ctx context.Context, user_id bson.ObjectID, refresh_token string) error
	Insert(ctx context.Context, user models.User) (bson.ObjectID, error)
}

type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Insert(ctx context.Context, user models.User) (bson.ObjectID, error) {
	collection := r.db.Database(db.Name).Collection(models.UserCollection)

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		var writeException mongo.WriteException
		if errors.As(err, &writeException) {
			for _, we := range writeException.WriteErrors {
				if we.Code == 11000 {
					switch {
					case strings.Contains(we.Message, "username_1"):
						return bson.NilObjectID, fmt.Errorf("username is already in use")
					case strings.Contains(we.Message, "email_1"):
						return bson.NilObjectID, fmt.Errorf("email is already in use")
					case strings.Contains(we.Message, "institute_email_1"):
						return bson.NilObjectID, fmt.Errorf("institute email is already in use")
					default:
						return bson.NilObjectID, fmt.Errorf("duplicate key error")
					}
				}
			}
		}
		return bson.NilObjectID, fmt.Errorf("failed to insert user: %w", err)
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, fmt.Errorf("inserted id is not ObjectID")
	}
	return id, nil
}
func (r *repository) InsertRefreshToken(
	ctx context.Context,
	user_id bson.ObjectID,
	refresh_token string,
) error {
	collection := r.db.Database(db.Name).Collection(models.TokenCollection)

	result, err := collection.InsertOne(ctx, models.RefreshToken{
		UserId:       user_id,
		RefreshToken: refresh_token,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return err
	}
	return nil
}
