package models


type Client struct{
	ID string `bson:"id" json:"id"`
	User string `bson:"user" json:"user"`
	// Conn *websocket.com
	Send chan[]byte
}