package peer

import (
	"context"
	"ep-peer-service/internal/db"
	"ep-peer-service/internal/features/common"
	"ep-peer-service/internal/models"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	GetPeer(
		ctx context.Context,
		user_id bson.ObjectID,
	) (*models.Peer, error)
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
