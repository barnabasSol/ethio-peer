package sessions

import (
	"context"
	"ep-streaming-service/internal/db"
	"ep-streaming-service/internal/features/common/pagination"
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
	GetSessions(
		context.Context,
		pagination.Pagination,
		string,
	) (*[]Session, error)
	IsOwner(
		ctx context.Context,
		session_id, username string,
	) (bool, error)
	UpdateSession(context.Context, Update) error
	InsertSession(
		ctx context.Context,
		session Create,
		username string,
		user_id string,
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
			"invalid session id",
		)
	}
	var sess models.Session
	err = collection.FindOne(ctx, bson.M{
		"_id":            oid,
		"owner.username": username,
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
	oid, err := bson.ObjectIDFromHex(req.SessionId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid user id",
		)
	}
	filter := bson.M{"_id": oid}
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
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed updating session",
		)
	}
	return nil
}

func (r *repository) GetSessions(
	ctx context.Context,
	p pagination.Pagination,
	req string,
) (*[]Session, error) {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)

	projection := bson.D{
		{Key: "_id", Value: 1},
		{Key: "session_name", Value: 1},
		{Key: "owner", Value: 1},
		{Key: "description", Value: 1},
		{Key: "starts_at", Value: 1},
		{Key: "ended_at", Value: 1},
		{Key: "participants", Value: 1},
	}

	popts := p.GetOptions().SetProjection(projection)

	var filter bson.D
	now := time.Now().UTC()
	switch req {
	case "upcoming":
		filter = bson.D{{
			Key: "starts_at",
			Value: bson.D{{
				Key:   "$gte",
				Value: now,
			}},
		}}
		popts.SetSort(bson.D{{
			Key:   "starts_at",
			Value: -1,
		}})
	case "ongoing":
		filter = bson.D{{
			Key: "starts_at",
			Value: bson.D{{
				Key:   "$lte",
				Value: now,
			}},
		},
			{
				Key: "ended_at",
				Value: bson.D{{
					Key:   "$in",
					Value: []any{nil},
				}},
			}}
		popts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	case "concluded":
		filter = bson.D{
			{
				Key: "ended_at",
				Value: bson.D{{
					Key:   "$ne",
					Value: nil,
				}},
			}}
		popts.SetSort(bson.D{{Key: "ended_at", Value: -1}})
	default:
		filter = bson.D{}
	}

	cursor, err := collection.Find(ctx, filter, popts)

	if err != nil {
		log.Println(err)
		return &[]Session{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed fetching sessions",
		)
	}

	defer cursor.Close(ctx)

	var result []Session
	if err := cursor.All(ctx, &result); err != nil {
		log.Println(err)
		return &[]Session{}, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed mapping session response",
		)
	}
	if result == nil {
		return &[]Session{}, nil
	}
	return &result, nil
}

func (r *repository) InsertSession(
	ctx context.Context,
	session Create,
	username string,
	user_id string,
) (string, error) {

	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	starts_at := time.Now().UTC()
	if session.StartsAt != nil {
		starts_at = *session.StartsAt
	}
	result, err := collection.InsertOne(ctx, models.Session{
		SessionName:   session.Name,
		Description:   session.Description,
		Tags:          session.Tags,
		Topic:         models.Topic(session.Topic),
		Scores:        []models.SessionScore{},
		ComputedScore: "0",
		Participants:  []models.Participant{},
		Owner: models.Owner{
			UserId:         user_id,
			Username:       username,
			Name:           session.OwnerName,
			ProfilePicture: session.OwnerProfilePic,
		},
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
		StartsAt:  &starts_at,
		EndedAt:   nil,
	})
	if err != nil {
		log.Println(err)
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
