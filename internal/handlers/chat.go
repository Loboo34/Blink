package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Loboo34/Blink/internal/database"
	"github.com/Loboo34/Blink/internal/models"
	"github.com/Loboo34/Blink/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed", nil)
		return
	}

	// Default limit: 50 messages, max 100
	limit := int64(50)
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 64); err == nil && parsedLimit > 0 {
			if parsedLimit > 100 {
				limit = 100
			} else {
				limit = parsedLimit
			}
		}
	}

	collection := database.DB.Collection("messages")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query latest messages sorted descending by timeStamp
	findOptions := options.Find().
		SetSort(bson.D{{Key: "timeStamp", Value: -1}}).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch messages", nil)
		return
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error reading messages", nil)
		return
	}

	// If no messages found, return an empty slice rather than null in JSON
	if messages == nil {
		messages = []models.Message{}
	}

	// Reverse slice to return messages in chronological order (oldest to newest)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	utils.RespondWithJson(w, http.StatusOK, "Messages retrieved successfully", messages)
}
