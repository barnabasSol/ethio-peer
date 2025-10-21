package analytics

import (
	"context"
	"ep-streaming-service/internal/db"
	"ep-streaming-service/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	GetDailyAnalyticsAggregation(ctx context.Context) (*SessionAnalytics, error)
	GetHourlyAnalyticsAggregation(ctx context.Context) (*SessionAnalytics, error)
}
type repository struct {
	ses_col *mongo.Collection
	db      *mongo.Client
}

func NewRepository(m *mongo.Client) Repository {
	return &repository{
		db:      m,
		ses_col: m.Database(db.Name).Collection(models.SessionCollection),
	}
}

func (r *repository) GetDailyAnalyticsAggregation(
	ctx context.Context,
) (*SessionAnalytics, error) {

	pipeline := GetDailyAnalyticsPipeline()
	cur, err := r.ses_col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to fetch aggregate report, try again later",
		)
	}
	defer cur.Close(ctx)

	var w []Weekly
	if err := cur.All(ctx, &w); err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"transformation failure",
		)
	}
	return &SessionAnalytics{
		Weekly: &w,
	}, nil
}

func (r *repository) GetHourlyAnalyticsAggregation(
	ctx context.Context,
) (*SessionAnalytics, error) {

	pipeline := GetHourlyAnalyticsPipeline()
	cur, err := r.ses_col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to fetch aggregate report, try again later",
		)
	}
	defer cur.Close(ctx)

	var h []Hourly
	if err := cur.All(ctx, &h); err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"transformation failure",
		)
	}
	return &SessionAnalytics{
		Hourly: &h,
	}, nil
}
