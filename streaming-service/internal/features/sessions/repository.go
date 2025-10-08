package sessions

import (
	"context"
	"ep-streaming-service/internal/db"
	"ep-streaming-service/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	DeleteSession(ctx context.Context)
	FilterSessions(context.Context, string)
	IsOwner(ctx context.Context, session_id, username string) (bool, error)
	UpdateSession(context.Context, Update) error
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

func (r *repository) IsOwner(
	ctx context.Context,
	session_id string,
	username string,
) (bool, error) {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	oid, err := bson.ObjectIDFromHex(session_id)
	if err != nil {
		return false, echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid user id",
		)
	}
	var sess models.Session
	err = collection.FindOne(ctx, bson.M{
		"owner.username": username,
		"_id":            oid,
	}).Decode(&sess)

	if err != nil {
		return false, echo.NewHTTPError(
			http.StatusNotFound,
			`user not found`,
		)
	}

	return sess.Owner.Username == username, nil
}

func (r *repository) UpdateSession(ctx context.Context, req Update) error {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	filter := bson.M{"_id": req.SessionId}
	update_filter := bson.M{}
	if req.SessionName != nil && *req.SessionName != "" {
		update_filter["session_name"] = *req.SessionName
	}
	if req.Description != nil && *req.Description != "" {
		update_filter["description"] = *req.Description
	}
	if req.IsEnded != nil && *req.IsEnded {
		update_filter["ended_at"] = time.Now().UTC()
	}

	update := bson.M{"$set": update_filter}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed updating session",
		)
	}
	return nil
}

func (r *repository) FilterSessions(context.Context, string) {
	panic("unimplemented")
}
func (r *repository) InsertSession(
	ctx context.Context,
	session Create,
	username string,
) (string, error) {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	result, err := collection.InsertOne(ctx, models.Session{
		SessionName:  session.Name,
		Description:  session.Description,
		Tags:         session.Tags,
		Topic:        models.Topic(session.Topic),
		Participants: []models.Participant{},
		Owner: models.Owner{
			Username:       username,
			Name:           session.OwnerName,
			ProfilePicture: session.OwnerProfilePic,
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
