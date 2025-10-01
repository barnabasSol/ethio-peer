package otp

import (
	"context"
	"time"
)

func scan(ctx context.Context, m *OTPManager) {
	ticker := time.NewTicker(5 * time.Second)
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
	// log.Println("cleaning")
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for k, otp := range m.collection {
		if otp.TTL.Before(now) {
			// log.Println("cleaned")
			delete(m.pending_users, otp.UserId)
			delete(m.collection, k)
		}
	}

}
