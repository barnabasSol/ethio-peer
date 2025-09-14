package main

import (
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
}

// func generateToken(room, identity string) (string, error) {

// 	canPublish := true
// 	canSubscribe := true

// 	at := auth.NewAccessToken(apiKey, apiSecret)
// 	grant := &auth.VideoGrant{
// 		RoomJoin:     true,
// 		Room:         room,
// 		CanPublish:   &canPublish,
// 		CanSubscribe: &canSubscribe,
// 	}
// 	at.SetVideoGrant(grant).
// 		SetIdentity(identity).
// 		SetValidFor(time.Hour)

// 	return at.ToJWT()
// }
