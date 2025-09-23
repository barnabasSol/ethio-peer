package resetpassword

import (
	"ep-auth-service/internal/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	UpdatePassword(user_id, password_hash string) error
	GetUser(VerifyRequest) (*models.User, error)
}

type repository struct {
	db *mongo.Client
}
