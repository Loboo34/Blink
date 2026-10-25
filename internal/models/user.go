package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName  string             `bson:"fullname" json:"fullname"`
	Password  string             `bson:"password,omitempty" json:"-"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
