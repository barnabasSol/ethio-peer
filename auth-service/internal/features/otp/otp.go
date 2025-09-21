package otp

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OTP struct {
	UserId    string    `json:"user_id"`
	SessionId string    `json:"session_id"`
	Value     string    `json:"value"`
	TTL       time.Time `json:"ttl"`
}

type OTPManager struct {
	mu            sync.RWMutex
	collection    map[string]OTP
	pending_users map[string]struct{}
}

func NewOTPManager(ctx context.Context) *OTPManager {
	m := &OTPManager{
		mu:            sync.RWMutex{},
		collection:    make(map[string]OTP),
		pending_users: map[string]struct{}{},
	}
	go scan(ctx, m)
	return m
}

func (m *OTPManager) removeOTP(session_id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.collection, session_id)
}

func (m *OTPManager) Generate(user_id string) (*OTP, error) {
	if m.exists(user_id) {
		return nil, echo.NewHTTPError(
			http.StatusConflict,
			"otp already pending, please try again later",
		)
	}
	m.add_to_pending(user_id)
	exp := os.Getenv("OTP_EXP_IN_MINS")
	expInMins, err := strconv.Atoi(exp)
	if err != nil {
		log.Println("Invalid OTP_EXP_IN_MINS:", err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate otp",
		)
	}

	value, err := generateOTP()

	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate otp",
		)
	}

	sessionID := uuid.NewString()
	otp := OTP{
		UserId:    user_id,
		SessionId: sessionID,
		Value:     value,
		TTL:       time.Now().Add(time.Duration(expInMins) * time.Minute),
	}

	m.mu.Lock()
	m.collection[sessionID] = otp
	m.mu.Unlock()
	return &otp, nil
}

func generateOTP() (string, error) {
	const otpLength = 6
	const digits = "0123456789"

	otp := make([]byte, otpLength)
	for i := range otpLength {
		num, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(digits))),
		)
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}

func (m *OTPManager) exists(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.pending_users[userID]
	return exists
}

func (m *OTPManager) add_to_pending(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pending_users[userID] = struct{}{}
}
