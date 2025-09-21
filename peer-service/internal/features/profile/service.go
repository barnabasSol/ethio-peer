package profile

import (
	"context"
	"ep-peer-service/internal/features/common"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service interface {
	UpdateProfilePicture(
		ctx context.Context,
		user_id string,
	) (string, error)

	DeleteProfilePicture(
		ctx context.Context,
		user_id string,
	) error
}

type service struct {
	r         Repository
	minio     *minio.Client
	minio_cfg common.MinioConfig
}

func NewService(m *minio.Client, repo Repository) Service {
	return &service{
		minio:     m,
		r:         repo,
		minio_cfg: common.GetMinioConfig(),
	}
}

func (s *service) UpdateProfilePicture(
	ctx context.Context,
	user_id string,
) (string, error) {
	object_key := uuid.NewString()
	expiry := time.Minute * 15

	presignedURL, err := s.minio.PresignedPutObject(
		ctx,
		s.minio_cfg.ImageBucket,
		object_key,
		expiry,
	)
	if err != nil {
		return "", err
	}
	id, err := bson.ObjectIDFromHex(user_id)
	if err != nil {
		return "", common.ErrInvalidPeerId
	}
	p, err := s.r.GetPeer(ctx, id)
	if err != nil {
		return "", err
	}
	if p.UserId != id {
		return "", echo.NewHTTPError(
			http.StatusForbidden,
			"cannot change profile picture besides your own",
		)
	}

	s.r.UpdateProfilePicture(
		ctx,
		id,
		fmt.Sprintf(
			"https://%s/%s/%s",
			s.minio_cfg.URL,
			s.minio_cfg.ImageBucket,
			object_key,
		),
	)
	return presignedURL.String(), nil

}

func (s *service) DeleteProfilePicture(
	ctx context.Context,
	userID string,
) error {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid peer ID",
		)
	}
	p, err := s.r.GetPeer(ctx, id)
	if err != nil {
		return err
	}

	if p.UserId != id {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"cannot delete profile picture besides your own",
		)
	}
	if p.ProfilePhoto == "" {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"no profile picture to delete",
		)
	}
	object_key := strings.TrimPrefix(
		p.ProfilePhoto,
		fmt.Sprintf(
			"https://%s/%s/",
			s.minio_cfg.URL,
			s.minio_cfg.ImageBucket,
		),
	)
	err = s.minio.RemoveObject(
		ctx,
		s.minio_cfg.ImageBucket,
		object_key,
		minio.RemoveObjectOptions{},
	)

	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to delete profile picture",
		)
	}

	if err := s.r.UpdateProfilePicture(ctx, id, ""); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to update peer profile",
		)
	}

	return nil
}
