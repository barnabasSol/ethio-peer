package scoring

import (
	"context"
	"ep-streaming-service/internal/db"
	"ep-streaming-service/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	InsertScore(context.Context, Score, string) error
}

type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) InsertScore(
	ctx context.Context,
	req Score,
	user_id string,
) error {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	ssid, err := bson.ObjectIDFromHex(req.SessionId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid session id",
		)
	}

	ss := models.SessionScore{
		UserId:    user_id,
		Score:     req.Score,
		Comment:   req.Comment,
		CreatedAt: time.Now().UTC(),
	}

	filter := bson.D{
		{Key: "_id", Value: ssid},
	}
	update := bson.M{
		"$push": bson.M{"scores": ss},
		"$set":  bson.M{"updated_at": time.Now().UTC()},
	}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to store score",
		)
	}
	return nil
}
