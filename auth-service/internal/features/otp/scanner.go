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
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for k, v := range m.collection {
		if v.TTL.Before(now) {
			delete(m.collection, k)
		}
	}

}
