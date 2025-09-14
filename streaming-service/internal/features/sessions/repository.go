package sessions

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	GetLiveSessions()
	DeleteSession()
	CreateSession()
}

type repository struct {
	db *mongo.Client
}
