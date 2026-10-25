package websocket

import (
	"net/http"
	"os"
	"strings"

	"github.com/Loboo34/Blink/internal/auth"
	"github.com/Loboo34/Blink/internal/utils"
	"github.com/gorilla/websocket"
)


func allowedOrigins() map[string]bool {
	origins := make(map[string]bool)
	raw := os.Getenv("ALLOWED_ORIGINS")
	if raw != "" {
		for _, o := range strings.Split(raw, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				origins[o] = true
			}
		}
	}

	origins["http://localhost:3000"] = true
	origins["http://localhost:8080"] = true
	origins["http://127.0.0.1:3000"] = true
	origins["http://127.0.0.1:8080"] = true
	return origins
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		
		if origin == "" {
			return true
		}
		return allowedOrigins()[origin]
	},
}

func ServeWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			tokenString = r.Header.Get("Authorization")
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		if tokenString == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing authentication token", nil)
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token: "+err.Error(), nil)
			return
		}

		userID, ok := claims["id"].(string)
		if !ok || userID == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID in token", nil)
			return
		}

		userName, _ := claims["name"].(string)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			utils.Logger.Warn("WebSocket upgrade failed: " + err.Error())
			return
		}

		client := &Client{
			hub:      hub,
			User:     userID,
			UserName: userName,
			Conn:     conn,
			Send:     make(chan Message, 250),
		}

		hub.register <- client

		go client.writePump()
		go client.readPump()
	}
}
