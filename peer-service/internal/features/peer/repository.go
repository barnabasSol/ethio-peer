package peer

import (
	"context"
	"ep-peer-service/internal/db"
	"ep-peer-service/internal/features/common"
	"ep-peer-service/internal/models"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository interface {
	GetPeer(
		ctx context.Context,
		user_id bson.ObjectID,
	) (*models.Peer, error)
	GetTopPeers(
		ctx context.Context,
	) (*[]TopPeer, error)
}

type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{db}
}

func (r *repository) GetPeer(
	ctx context.Context,
	user_id bson.ObjectID,
) (*models.Peer, error) {
	peer_collection := r.db.Database(db.Name).Collection(models.PeerCollection)
	var peer models.Peer

	filter := bson.M{"_id": user_id}

	err := peer_collection.FindOne(ctx, filter).Decode(&peer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, common.ErrPeerNotFound
		}
		log.Println(err)
		return nil, err
	}
	return &peer, nil

}

func (r *repository) GetTopPeers(
	ctx context.Context,
) (*[]TopPeer, error) {
	peer_collection := r.db.Database(db.Name).Collection(models.PeerCollection)
	var top_peers []TopPeer

	find_options := options.Find().
		SetSort(bson.D{{Key: "overall_score", Value: -1}}).
		SetLimit(3).
		SetProjection(bson.D{
			{Key: "_id", Value: 1},
			{Key: "overall_score", Value: 1},
			{Key: "profile_photo", Value: 1},
		})

	cursor, err := peer_collection.Find(ctx, bson.M{}, find_options)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, common.ErrPeerNotFound
		}
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var top_peer TopPeer
		if err := cursor.Decode(&top_peer); err != nil {
			log.Println(err)
			return nil, echo.NewHTTPError(
				http.StatusInternalServerError,
				"transformation failure",
			)
		}
		top_peer.Id = top_peer.OID.Hex()
		top_peers = append(top_peers, top_peer)
	}

	if err := cursor.Err(); err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed mapping top peers",
		)
	}
	if len(top_peers) == 0 {
		return nil, echo.NewHTTPError(
			http.StatusNotFound,
			"no top peers yet",
		)
	}

	return &top_peers, nil
}
