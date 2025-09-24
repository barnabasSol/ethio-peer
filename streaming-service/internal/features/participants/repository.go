package participants

import (
	"context"
	"ep-streaming-service/internal/db"
	"ep-streaming-service/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	UpdateFlag(ctx context.Context, flag Flag) error
	Insert(context.Context, Join) error
	GetSession(context.Context, string) (*models.Session, error)
}

type repository struct {
	db *mongo.Client
}

func (r *repository) GetSession(
	ctx context.Context,
	sid string,
) (*models.Session, error) {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	var session models.Session
	err := collection.FindOne(
		ctx,
		bson.M{"_id": sid},
	).Decode(&session)

	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusNotFound,
			"session not found",
		)
	}
	if session.EndedAt != nil {
		return nil, echo.NewHTTPError(
			http.StatusNotFound,
			"session expired",
		)
	}

	return &session, nil
}

// Insert implements Repository.
func (r *repository) Insert(context.Context, Join) error {
	panic("unimplemented")
}

// UpdateFlag implements Repository.
func (r *repository) UpdateFlag(
	ctx context.Context,
	flag Flag,
) error {
	panic("unimplemented")
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

// user_collection := r.db.Database(db.Name).Collection(models.UserCollection)
// filter := bson.M{"_id": user_id}
// var user models.User

// err := user_collection.FindOne(ctx, filter).Decode(&user)
// if err != nil {
// 	return nil, errors.New("user not found")
// }
// return &user, nil
