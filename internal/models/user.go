package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"id,omitempty" json:"id"`
	FullName  string             `bson:"fullname" json:"fullname"`
	Password  string             `bson:"password,omitempty" json:"password"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
