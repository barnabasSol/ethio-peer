package resetpassword

import (
	"context"
	"ep-auth-service/internal/db"
	"ep-auth-service/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	UpdatePassword(
		ctx context.Context,
		institute_email string,
		password_hash string,
	) error
}

type repository struct {
	db *mongo.Client
}

func NewRepository(
	db *mongo.Client,
) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) UpdatePassword(
	ctx context.Context,
	institute_email string,
	password_hash string,
) error {

	user_collection := r.db.Database(db.Name).Collection(models.UserCollection)
	filter := bson.D{{Key: "institute_email", Value: institute_email}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "password_hash", Value: password_hash},
			{Key: "updated_at", Value: time.Now().UTC()},
		}},
	}

	result, err := user_collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to update profile picture",
		)
	}

	if result.MatchedCount == 0 {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"user not found",
		)
	}

	return nil
}
