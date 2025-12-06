package hamdlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Loboo34/Blink/internal/auth"
	"github.com/Loboo34/Blink/internal/database"
	"github.com/Loboo34/Blink/internal/models"
	"github.com/Loboo34/Blink/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func Register(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST Allowed", "")
		return
	}

	var req struct{
		FullName string `json:"fullname"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req);err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid json format", "")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error hashing password", "")
		return
	}

	collection := database.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser := models.User{
		ID : primitive.NewObjectID(),
		FullName: req.FullName,
		Password: hashedPassword,
		Role: "Member",
		CreatedAt: time.Now(),
	}

	_,err = collection.InsertOne(ctx, newUser)
	if err != nil{
		utils.Logger.Warn("Failed to Register user")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error registering User", "")
		return
	}

	utils.Logger.Info("Registration Successful")
	utils.RespondWithJson(w, http.StatusCreated, "User created", map[string]interface{}{
		"Name": req.FullName,
	})	
}

func Login(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Post Allowed", "")
		return
	}

	var req struct{
		FullName string `json:"fullname"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req);err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid json format", "")
		return
	}

	collection := database.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"fullname": req.FullName}).Decode(&user)
	if err != nil{
		if err == mongo.ErrNoDocuments{
			utils.RespondWithError(w, http.StatusNotFound, "User not found", "")
		}else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding user", "")
		}
		return
	}

	if !auth.ComparePassword(req.Password, user.Password){
		utils.RespondWithError(w, http.StatusInternalServerError, "Invalid Password", "")
		return
	}

	token, err := auth.GenerateJWTToken(user.ID.Hex(), user.FullName, user.Role)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating Token", "")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Login Successfull", map[string]interface{}{"token": token, "user": map[string]interface{}{
		"id":       user.ID.Hex(),
		"fullname": user.FullName,
	}})
}