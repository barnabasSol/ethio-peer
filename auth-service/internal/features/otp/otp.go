package otp

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	SessionId string    `json:"session_id"`
	Value     string    `json:"value"`
	TTL       time.Time `json:"ttl"`
}

type OTPManager struct {
	mu         sync.RWMutex
	collection map[string]OTP
}

func NewOTPManager(ctx context.Context) *OTPManager {
	m := &OTPManager{
		mu:         sync.RWMutex{},
		collection: make(map[string]OTP),
	}
	go scan(ctx, m)
	return m
}

func (m *OTPManager) Generate() (*OTP, error) {
	exp := os.Getenv("OTP_EXP_IN_MINS")
	expInMins, err := strconv.Atoi(exp)
	if err != nil {
		log.Fatal("Invalid OTP_EXP_IN_MINS:", err)
		return nil, err
	}

	value, err := generateOTP()

	if err != nil {
		return nil, err
	}

	sessionID := uuid.NewString()
	otp := OTP{
		SessionId: sessionID,
		Value:     value,
		TTL:       time.Now().Add(time.Duration(expInMins) * time.Minute),
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.collection[sessionID] = otp
	return &otp, nil
}

func generateOTP() (string, error) {
	const otpLength = 6
	const digits = "0123456789"

	otp := make([]byte, otpLength)
	for i := range otpLength {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}
