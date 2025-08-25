package otp

import (
	"context"
	"time"
)

func scan(ctx context.Context, m *OTPManager) {
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanup(m)
		case <-ctx.Done():
			return
		}
	}
}

func cleanup(m *OTPManager) {

}
