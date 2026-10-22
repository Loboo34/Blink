package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Loboo34/Blink/internal/database"
	"github.com/Loboo34/Blink/internal/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// load env
func init() {
	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load()
	}
	if err != nil {
		log.Fatal("Failed to load .env")
	}
}

func main() {

	r := mux.NewRouter()
	r.Use()

	utils.InitLogger(os.Getenv("ENV") == "production")
	defer utils.Logger.Sync()

	db := database.ConnectDB()
	fmt.Println("connected to: ", db.Name())
	fmt.Println("Hello")

	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.Logger.Info("server running on  port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
