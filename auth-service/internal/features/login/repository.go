package login

import (
	"context"
	"ep-auth-service/internal/db"
	"ep-auth-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	GetUser(ctx context.Context, login LoginRequest) (*models.User, error)
	InsertRefreshToken(ctx context.Context, user_id bson.ObjectID, refresh_token string) error
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
	login LoginRequest,
) (*models.User, error) {
	user_collection := r.db.Database(db.Name).Collection(models.UserCollection)
	filter := bson.D{}
	if login.Username != nil {
		filter = bson.D{{Key: "username", Value: *login.Username}}
	} else if login.Email != nil {
		filter = bson.D{{Key: "email", Value: *login.Email}}
	} else if login.InstituteEmail != nil {
		filter = bson.D{{Key: "institute_email", Value: *login.InstituteEmail}}
	}
	var user models.User
	err := user_collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, ErrUserDoesntExist
	}
	return &user, nil
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
