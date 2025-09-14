package profile

import (
	"context"
	"ep-peer-service/internal/db"
	"ep-peer-service/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	GetPeer(
		ctx context.Context,
		user_id bson.ObjectID,
	) (*models.Peer, error)
	UpdateProfilePicture(
		ctx context.Context,
		user_id bson.ObjectID,
		profile_photo_url string,
	) error
}
type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) UpdateProfilePicture(
	ctx context.Context,
	userID bson.ObjectID,
	profilePhotoURL string,
) error {
	peerCollection := r.db.Database(db.Name).Collection(models.PeerCollection)

	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "profile_photo", Value: profilePhotoURL},
			{Key: "updated_at", Value: time.Now().UTC()},
		}},
	}

	result, err := peerCollection.UpdateOne(ctx, filter, update)
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

func (r *repository) GetPeer(
	ctx context.Context,
	user_id bson.ObjectID,
) (*models.Peer, error) {
	peer_collection := r.db.Database(db.Name).Collection(models.PeerCollection)
	filter := bson.D{{Key: "_id", Value: user_id}}
	var peer models.Peer
	err := peer_collection.FindOne(ctx, filter).Decode(&peer)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to fetch peer",
		)

	}
	return &peer, nil
}
