package analytics

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Service interface {
	GetSessionAnalytics(ctx context.Context, filter string) (*SessionAnalytics, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetSessionAnalytics(
	ctx context.Context,
	filter string,
) (*SessionAnalytics, error) {
	switch filter {
	case "weekly":
		return s.repo.GetDailyAnalyticsAggregation(ctx)
	case "yearly":
		return s.repo.GetMonthlyAnalyticsAggregation(ctx)
	default:
		return nil, echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid filter",
		)
	}
}
