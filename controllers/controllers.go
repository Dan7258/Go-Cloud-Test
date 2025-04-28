// Package controllers предоставляет обработчики HTTP-запросов для управления клиентами в системе.
package controllers

import (
	"cloud/logger"
	"cloud/models"
	"encoding/json"
	"net/http"
	"strings"
)

// ClientHandler обрабатывает HTTP-запросы для работы с клиентами.
// В зависимости от метода запроса вызывает соответствующую функцию.
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

// CreateClient создает нового клиента на основе данных, переданных в теле запроса.
// В случае ошибки отправляет соответствующий HTTP-статус и сообщение.
func CreateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	client := new(models.RateLimits)
	json.NewDecoder(r.Body).Decode(client)

	if client.ClientID == "" {
		logger.SendError(w, http.StatusInternalServerError, "Запрос не содержит данные")
		return
	}

	err := models.CreateClient(*client)
	if err != nil {
		logger.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetClient возвращает информацию о клиенте по его идентификатору (clientID).
// Если идентификатор отсутствует или происходит ошибка — отправляет сообщение об ошибке.
func GetClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	splitPath := strings.Split(r.URL.Path, "/")
	if len(splitPath) < 3 {
		logger.SendError(w, http.StatusNotFound, "Отсутствует clientID")
		return
	}

	clientID := splitPath[2]
	client, err := models.GetClient(clientID)
	if err != nil {
		logger.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(client)
}

// UpdateClient обновляет информацию о клиенте на основе данных из тела запроса.
// При возникновении ошибки возвращает соответствующее сообщение об ошибке.
func UpdateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	updClient := new(models.RateLimits)
	json.NewDecoder(r.Body).Decode(updClient)

	err := models.UpdateClient(updClient)
	if err != nil {
		logger.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteClient удаляет клиента по идентификатору (clientID), переданному в URL.
// В случае ошибки отправляет сообщение об ошибке.
func DeleteClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	splitPath := strings.Split(r.URL.Path, "/")
	if len(splitPath) < 3 {
		logger.SendError(w, http.StatusNotFound, "Отсутствует clientID")
		return
	}

	clientID := splitPath[2]
	err := models.DeleteClient(clientID)
	if err != nil {
		logger.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
