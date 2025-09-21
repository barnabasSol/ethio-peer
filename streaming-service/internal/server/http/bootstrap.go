package server

import "os"

func (s *Server) bootstrap() error {
	lk_api_key := os.Getenv("LK_API_KEY")
	lk_api_secret := os.Getenv("LK_API_SECRET")
	_ = lk_api_key
	_ = lk_api_secret

	lk_egress_key := os.Getenv("LK_EGRESS_KEY")
	lk_egress_secret := os.Getenv("LK_EGRESS_SECRET")
	_ = lk_egress_key
	_ = lk_egress_secret
	return nil
}
