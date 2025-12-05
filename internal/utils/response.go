package utils

import (
	"encoding/json"
	"net/http"
)


type APIResponse struct{
Success bool `json:"success"`
Error string `json:"error"`
Message string `json:"message"`
Data interface{} `json:"data"`

}

func RespondWithError(w http.ResponseWriter, code int, message string, payload interface{}){
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)

	apiResponse := APIResponse{
		Success: false,
		Error: message,
		Data: payload,
	}

	json.NewEncoder(w).Encode(apiResponse)
}

func RespondWithJson(w http.ResponseWriter, code int, message string, payload interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	apiResponse :=APIResponse{
		Success: true,
		Message: message,
		Data: payload,
	}

	json.NewEncoder(w).Encode(apiResponse)
}