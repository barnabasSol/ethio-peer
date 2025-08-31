package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/livekit/protocol/auth"
)

func main() {
	godotenv.Load()
	fmt.Println(generateToken("alone", "barney"))
}

func generateToken(room, identity string) (string, error) {
	apiKey := os.Getenv("LK_API_KEY")
	apiSecret := os.Getenv("LK_API_SECRET")

	canPublish := true
	canSubscribe := true

	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         room,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscribe,
	}
	at.SetVideoGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}
