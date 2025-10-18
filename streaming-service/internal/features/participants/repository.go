package participants

import (
	"context"
	"ep-streaming-service/internal/db"
	"ep-streaming-service/internal/features/common/flags"
	"ep-streaming-service/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository interface {
	UpdateFlag(ctx context.Context, flag Flag) error
	GetParticipantByUsername(ctx context.Context, username string) (*models.Participant, error)
	Insert(context.Context, bool, Join) error
	GetSession(context.Context, bson.ObjectID) (*models.Session, error)
}

type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetSession(
	ctx context.Context,
	sid bson.ObjectID,
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
			http.StatusConflict,
			"session is over",
		)
	}

	return &session, nil
}

func (r *repository) Insert(
	ctx context.Context,
	is_owner bool,
	join Join,
) error {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	sessionID, err := bson.ObjectIDFromHex(join.SessionId)
	if err != nil {
		return err
	}

	participant := models.Participant{
		Username:       join.Username,
		Name:           join.Name,
		UserId:         join.UserId,
		ProfilePicture: join.ProfilePicture,
		IsAnonymous:    join.AsAnonymous,
		IsMuted:        false,
		FlagStatus:     flags.OK,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	filter := bson.M{"_id": sessionID}
	update := bson.M{
		"$push": bson.M{"participants": participant},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *repository) UpdateFlag(
	ctx context.Context,
	flag Flag,
) error {
	panic("unimplemented")
}

func (r *repository) GetParticipantByUsername(
	ctx context.Context,
	username string,
) (*models.Participant, error) {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	filter := bson.M{
		"participants.username": username,
	}
	projection := bson.M{
		"participants": bson.M{
			"$elemMatch": bson.M{
				"username": username,
			},
		},
	}

	var result struct {
		Participants []models.Participant `bson:"participants"`
	}
	err := collection.FindOne(
		ctx,
		filter,
		options.FindOne().SetProjection(projection),
	).Decode(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(result.Participants) == 0 {
		return nil, echo.NewHTTPError(
			http.StatusNotFound,
			"participant not found",
		)
	}
	return &result.Participants[0], nil
}
