package websoket

import "github.com/gorilla/websocket"

type Client struct {
	hub  *Hub
	User string `bson:"user" json:"user"`
	Conn *websocket.Conn
	Send chan Message
}
