package scoring

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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository interface {
	GetOwnerId(context.Context, bson.ObjectID) (string, error)
	InsertScore(context.Context, Score, string, string) error
	UpdateScore(context.Context, Score, string) error
	UpdateComputedScore(context.Context, bson.ObjectID, string) error
	GetAverageSessionScore(context.Context, string) (*float64, error)
	GetScores(context.Context, bson.ObjectID) ([]float32, error)
	IsScored(context.Context, bson.ObjectID, string) (bool, error)
}

type repository struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) UpdateComputedScore(
	ctx context.Context,
	oid bson.ObjectID,
	score string,
) error {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "computed_score", Value: score},
		}},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (r *repository) GetAverageSessionScore(
	ctx context.Context,
	oid string,
) (*float64, error) {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"owner.user_id":  oid,
			"computed_score": bson.M{"$ne": "0"},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": nil,
			"avgScore": bson.M{
				"$avg": bson.M{
					"$toDouble": "$computed_score",
				},
			},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		AvgScore float64 `bson:"avgScore"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, echo.NewHTTPError(
			http.StatusNotFound,
			"no sessions found for this user",
		)
	}
	return &result[0].AvgScore, nil
}

func (r *repository) GetOwnerId(
	ctx context.Context,
	sid bson.ObjectID,
) (string, error) {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	filter := bson.D{
		{Key: "_id", Value: sid},
	}
	projection := bson.D{
		{Key: "owner.user_id", Value: 1},
	}
	var result struct {
		Owner models.Owner `bson:"owner"`
	}
	err := collection.FindOne(
		ctx,
		filter,
		options.FindOne().SetProjection(projection),
	).Decode(&result)

	if err != nil {
		log.Println(err)
		if err == mongo.ErrNoDocuments {
			return "", echo.NewHTTPError(
				http.StatusNotFound,
				"session not found",
			)
		}
		return "", err
	}

	return result.Owner.UserId, nil
}

func (r *repository) GetScores(
	ctx context.Context,
	session_id bson.ObjectID,
) ([]float32, error) {

	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	filter := bson.D{
		{Key: "_id", Value: session_id},
	}
	projection := bson.D{
		{Key: "scores", Value: 1},
		{Key: "_id", Value: 0},
	}
	var result struct {
		Scores []models.SessionScore `bson:"scores"`
	}

	err := collection.FindOne(
		ctx,
		filter,
		options.FindOne().SetProjection(projection),
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	scores := make([]float32, len(result.Scores))
	for i, sc := range result.Scores {
		scores[i] = sc.Score
	}
	return scores, nil
}

func (r *repository) IsScored(
	ctx context.Context,
	session_id bson.ObjectID, user_id string,
) (bool, error) {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)

	count_filter := bson.D{
		{Key: "_id", Value: session_id},
		{Key: "scores.user_id", Value: user_id},
	}

	count, err := collection.CountDocuments(ctx, count_filter)
	if err != nil {
		return false, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed making rated check",
		)
	}
	return count > 0, nil

}

func (r *repository) UpdateScore(
	ctx context.Context,
	req Score,
	user_id string,
) error {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)

	soid, err := bson.ObjectIDFromHex(req.SessionId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid session id",
		)
	}
	filter := bson.D{
		{Key: "_id", Value: soid},
		{
			Key:   "scores.user_id",
			Value: user_id,
		},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "scores.$.score", Value: req.Score},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "scores.$.comment", Value: req.Comment},
		}},
	}

	_, err = collection.UpdateOne(
		ctx,
		filter,
		update,
	)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to update score",
		)
	}
	return nil
}

func (r *repository) InsertScore(
	ctx context.Context,
	req Score,
	username string,
	user_id string,
) error {
	collection := r.db.Database(db.Name).Collection(models.SessionCollection)
	ssid, err := bson.ObjectIDFromHex(req.SessionId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid session id",
		)
	}

	ss := models.SessionScore{
		UserId:    user_id,
		Username:  username,
		Score:     req.Score,
		Comment:   req.Comment,
		CreatedAt: time.Now().UTC(),
	}

	filter := bson.D{
		{Key: "_id", Value: ssid},
	}
	update := bson.M{
		"$push": bson.M{"scores": ss},
		"$set":  bson.M{"updated_at": time.Now().UTC()},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to store score",
		)
	}
	return nil
}
