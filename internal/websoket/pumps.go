package websoket




func (c *Client) readPump(){
	defer func ()  {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil{
			break
		}

		msg.SenderID = c.User
		c.hub.broadcast <- msg
	}

}


func (c *client) writePump(){
	defer c.Conn.Close()

	for msg := range c.send {
		 if err := c.conn.WriteJSON(msg); err != nil {
            break
        }
	}
}