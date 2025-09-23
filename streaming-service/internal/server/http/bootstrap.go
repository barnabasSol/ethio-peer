package server

import (
	"ep-streaming-service/internal/features/common/livekit"
	"ep-streaming-service/internal/features/sessions"
	"os"
)

func (s *Server) bootstrap() error {
	host := os.Getenv("LK_URL")
	lk_api_key := os.Getenv("LK_API_KEY")
	lk_api_secret := os.Getenv("LK_API_SECRET")
	lk_egress_key := os.Getenv("LK_EGRESS_KEY")
	lk_egress_secret := os.Getenv("LK_EGRESS_SECRET")
	livekit_cfg := livekit.NewConfig(
		host,
		lk_api_key,
		lk_api_secret,
		lk_egress_key,
		lk_egress_secret,
	)
	ss := sessions.NewService(nil, *livekit_cfg)

	return nil
}
