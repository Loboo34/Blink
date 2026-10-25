package websocket

import "github.com/gorilla/websocket"

type Client struct {
	hub      *Hub
	User     string `bson:"user" json:"user"`
	UserName string `bson:"userName" json:"userName"`
	Conn     *websocket.Conn
	Send     chan Message
}
