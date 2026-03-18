package web

import (
	"encoding/json"
	"net/http"
)

type SuccessResponseJson struct {
	Message string `json:"Message"`
}

type ErrorResponseJson struct {
	Error   string `json:"Error"`
	Message string `json:"Message"`
}

func SuccessResponse(mess string) interface{} {
	return SuccessResponseJson{
		Message: mess,
	}
}

func ErrorResponse(mess string) interface{} {
	return ErrorResponseJson{
		Error:   "ERROR",
		Message: mess,
	}
}

func UnauthorizedResponce() interface{} {
	return ErrorResponse("User unauthorized")
}

func SendJsonResponse(w http.ResponseWriter, status int, jsonMessage interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(jsonMessage)
}
