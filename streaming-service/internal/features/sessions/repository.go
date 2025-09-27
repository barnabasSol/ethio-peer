package sessions

import (
	"context"
	"ep-streaming-service/internal/db"
	"ep-streaming-service/internal/features/common/flags"
	"ep-streaming-service/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	DeleteSession(ctx context.Context)
	InsertSession(
		ctx context.Context,
		session Create,
		username string,
	) (string, error)
}

type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}
func (r *repository) InsertSession(
	ctx context.Context,
	session Create,
	username string,
) (string, error) {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	result, err := collection.InsertOne(ctx, models.Session{
		SessionName: session.Name,
		Description: session.Description,
		Tags:        session.Tags,
		Participants: []models.Participant{
			{
				Username:       username,
				Name:           session.OwnerName,
				ProfilePicture: session.OwnerProfilePic,
				IsOwner:        true,
				FlagStatus:     flags.OK,
				IsAnonymous:    false,
				CreatedAt:      time.Now().UTC(),
				UpdatedAt:      time.Now().UTC(),
			},
		},
		CreatedAt: time.Now().UTC(),
		EndedAt:   nil,
	})
	if err != nil {
		return "", echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to create session",
		)
	}
	if !result.Acknowledged {
		return "", echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to create session",
		)
	}
	id := result.InsertedID.(bson.ObjectID)
	return id.Hex(), nil

}

func (r *repository) DeleteSession(ctx context.Context) {
	panic("unimplemented")
}
