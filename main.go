package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Loboo34/Blink/internal/auth"
	"github.com/Loboo34/Blink/internal/database"
	handlers "github.com/Loboo34/Blink/internal/handlers"
	"github.com/Loboo34/Blink/internal/utils"
	"github.com/Loboo34/Blink/internal/websocket"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// load env
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env")
	}
}

func main() {
	r := mux.NewRouter()

	utils.InitLogger(os.Getenv("ENV") == "production")
	defer utils.Logger.Sync()

	if err := auth.InitJWT(); err != nil {
		utils.Logger.Fatal("Failed to initialize JWT:" + err.Error())
	}

	db, err := database.ConnectDB()
	if err != nil {
		utils.Logger.Fatal("Failed to connect to database: " + err.Error())
	}
	fmt.Println("connected to:", db.Name())

	// Init WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	//CORS 
	r.Use(utils.CORSMiddleware)

	// Authentication routes
	r.HandleFunc("/api/auth/register", handlers.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/login", handlers.Login).Methods("POST", "OPTIONS")

	//Chat history
	r.HandleFunc("/api/messages", auth.CheckAuth(handlers.GetMessages)).Methods("GET", "OPTIONS")

	// Chat WebSocket route
	r.HandleFunc("/ws", websocket.ServeWS(hub))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Logger.Info("server running on port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
