package websocket

import (
	"context"
	"strings"
	"time"

	"github.com/Loboo34/Blink/internal/database"
	"github.com/Loboo34/Blink/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		
		if strings.TrimSpace(msg.Content) == "" {
			continue
		}

	
		msg.ID = primitive.NewObjectID()
		msg.SenderID = c.User
		if c.UserName != "" {
			msg.SenderName = c.UserName
		}
		msg.TimeStamp = time.Now()

		
		go persistMessage(msg)

		
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}


func persistMessage(msg Message) {
	if database.DB == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.DB.Collection("messages").InsertOne(ctx, msg)
	if err != nil && utils.Logger != nil {
		utils.Logger.Warn("Failed to persist message to MongoDB: " + err.Error())
	}
}
