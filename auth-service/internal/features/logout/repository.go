package logout

import "go.mongodb.org/mongo-driver/v2/mongo"

type Repository interface {
	UpdateEmailVerified(user_id string) bool
}

type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) UpdateEmailVerified(user_id string) bool {
	return true
}
