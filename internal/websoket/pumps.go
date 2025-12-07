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


func (c *Client) writePump(){
	defer c.Conn.Close()

	for msg := range c.Send {
		 if err := c.Conn.WriteJSON(msg); err != nil {
            break
        }
	}
}