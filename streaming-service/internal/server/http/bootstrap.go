package server

import (
	broker "ep-streaming-service/internal/broker/rabbitmq"
	"ep-streaming-service/internal/features/common/livekit"
	"ep-streaming-service/internal/features/monitoring"
	"ep-streaming-service/internal/features/participants"
	"ep-streaming-service/internal/features/scoring"
	"ep-streaming-service/internal/features/sessions"
	"log"
	"os"
)

func (s *Server) bootstrap() error {
	host := os.Getenv("LK_URL")
	lk_api_key := os.Getenv("LK_API_KEY")
	lk_api_secret := os.Getenv("LK_API_SECRET")
	lk_egress_key := os.Getenv("LK_EGRESS_KEY")
	lk_egress_secret := os.Getenv("LK_EGRESS_SECRET")
	wh := os.Getenv("LK_WH_KEY")
	livekit_cfg := livekit.NewConfig(
		host,
		lk_api_key,
		lk_api_secret,
		lk_egress_key,
		lk_egress_secret,
		wh,
	)
	b, err := broker.InitRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	sr := sessions.NewRepository(s.db)
	ss := sessions.NewService(
		sr,
		b,
		*livekit_cfg,
	)
	sessions.InitHandler(
		ss,
		*livekit_cfg,
		s.echo.Group("session"),
	)

	pr := participants.NewRepository(s.db)
	ps := participants.NewService(pr, livekit_cfg)
	participants.InitHandler(
		ps,
		*livekit_cfg,
		s.echo.Group("participant"),
	)

	ms := monitoring.NewService(livekit_cfg, pr)
	monitoring.InitHandler(
		ms,
		*livekit_cfg,
		s.echo.Group("monitoring"),
	)

	scr := scoring.NewRepository(s.db)
	scs := scoring.NewService(scr, b)
	scoring.InitHandler(
		*livekit_cfg,
		s.echo.Group(""),
		scs,
	)
	return nil
}
