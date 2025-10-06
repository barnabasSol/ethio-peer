package refreshtoken

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
	GetUser(ctx context.Context, req Request) (*models.User, error)
	GetRefreshToken(ctx context.Context, req Request) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, req Request) error
	InsertRefreshToken(
		ctx context.Context,
		new_refresh_token string,
		req Request,
	) error
}

type repository struct {
	db *mongo.Client
}

func (r *repository) GetUser(ctx context.Context, req Request) (*models.User, error) {
	user_collection := r.db.Database(db.Name).Collection(models.UserCollection)
	user_obj_id, err := bson.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid user id",
		)
	}
	filter := bson.D{{Key: "_id", Value: user_obj_id}}
	var user models.User
	err = user_collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusNotFound,
			"user doesnt exist",
		)
	}
	return &user, nil
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetRefreshToken(
	ctx context.Context,
	req Request,
) (*models.RefreshToken, error) {

	token_collection := r.db.Database(db.Name).Collection(models.TokenCollection)
	userObjId, err := bson.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid user id",
		)
	}

	filter := bson.M{
		"user_id":       userObjId,
		"refresh_token": req.RefreshToken,
	}

	var refreshToken models.RefreshToken
	err = token_collection.FindOne(ctx, filter).Decode(&refreshToken)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusUnauthorized,
			"refresh token not found or invalid",
		)
	}

	return &refreshToken, nil

}

func (r *repository) DeleteRefreshToken(
	ctx context.Context,
	req Request,
) error {
	collection := r.db.Database(db.Name).Collection(models.TokenCollection)
	filter := bson.D{{Key: "user_id", Value: req.UserId}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to clear refresh token",
		)
	}
	return nil
}

func (r *repository) InsertRefreshToken(
	ctx context.Context,
	new_refresh_token string,
	req Request,
) error {
	collection := r.db.Database(db.Name).Collection(models.TokenCollection)

	user_obj_id, err := bson.ObjectIDFromHex(req.UserId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid user id",
		)
	}
	result, err := collection.InsertOne(ctx, models.RefreshToken{
		UserId:       user_obj_id,
		RefreshToken: new_refresh_token,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil || !result.Acknowledged {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"something went wrong during auth",
		)
	}
	return nil
}
