package changepassword

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	UpdatePassword(user_id string, password_hash string)
}

type repository struct {
	db *mongo.Client
}
