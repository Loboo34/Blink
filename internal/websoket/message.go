package websoket

import (
	"time"
)

type Message struct {
	SenderID   string    `bson:"senderID" json:"senderID"`
	SenderName string    `bson:"senderName" json:"senderName"`
	Content    string    `bson:"content" json:"content"`
	TimeStamp  time.Time `bson:"timeStamp" json:"timeStamp"`
}
