package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const UserCollection = "users"

type User struct {
	Id                     bson.ObjectID `bson:"_id,omitempty"`
	Username               string        `bson:"username"`
	Name                   string        `bson:"name"`
	PasswordHash           string        `bson:"password_hash"`
	Roles                  []string      `bson:"roles"`
	IsActive               bool          `bson:"is_active"`
	InstituteEmail         string        `bson:"institute_email"`
	InstituteEmailVerified bool          `bson:"institute_email_verified"`
	Email                  string        `bson:"email,omitempty"`
	CreatedAt              time.Time     `bson:"created_at"`
	UpdatedAt              time.Time     `bson:"updated_at"`
}
