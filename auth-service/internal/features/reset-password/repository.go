package resetpassword

import (
	broker "ep-auth-service/internal/broker/rabbitmq"
	"ep-auth-service/internal/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	UpdatePassword(user_id, password_hash string) error
	GetUser(VerifyRequest) (*models.User, error)
}

type repository struct {
	broker *broker.RabbitMQ
	db     *mongo.Client
}

func NewRepository(
	db *mongo.Client,
	br *broker.RabbitMQ,
) Repository {
	return &repository{
		db:     db,
		broker: br,
	}
}

func (r *repository) GetUser(req VerifyRequest) (*models.User, error) {
	panic("unimplemented")
}

func (r *repository) UpdatePassword(
	user_id string,
	password_hash string,
) error {
	panic("unimplemented")
}
