package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Message struct{
	ID primitive.ObjectID `bson:"id" json:"id"`
	SenderID string `bson:"senderID" json:"senderID"`
	SenderName string `bson:"senderName" json:"senderName"`
	Content string `bson:"content" json:"content"`
	TimeStamp time.Time `bson:"timeStamp" json:"timeStamp"`
}