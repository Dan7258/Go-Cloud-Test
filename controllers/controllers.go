package controllers

import (
	"cloud/models"
	"encoding/json"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetClient(w, r)
	case http.MethodPost:
		CreateClient(w, r)
	case http.MethodPatch:
		UpdateClient(w, r)
	case http.MethodDelete:
		DeleteClient(w, r)
	}
}

func CreateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	client := new(models.RateLimits)
	json.NewDecoder(r.Body).Decode(client)
	if client.ClientID == "" || client.Capacity == 0 || client.RatePerSecond == 0 {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Запрос не содержит данные",
		})
	}
	err := models.CreateClient(client)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	splitPath := strings.Split(r.URL.Path, "/")
	if len(splitPath) < 3 {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Отсутствует clientID",
		})
		return
	}

	clientID := splitPath[2]
	client, err := models.GetClient(clientID)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	updClient := new(models.RateLimits)
	json.NewDecoder(r.Body).Decode(updClient)
	err := models.UpdateClient(updClient)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	splitPath := strings.Split(r.URL.Path, "/")
	if len(splitPath) < 3 {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Отсутствует clientID",
		})
		return
	}
	clientID := splitPath[2]
	err := models.DeleteClient(clientID)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
