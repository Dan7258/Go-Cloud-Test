// Package logger предоставляет функции для логирования сообщений различных уровней
// и отправки ошибок в формате JSON через HTTP-ответ.
package logger

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse представляет структуру ответа об ошибке, отправляемого клиенту.
type ErrorResponse struct {
	Code    int    `json:"code"`    // HTTP-статус код ошибки.
	Message string `json:"message"` // Сообщение об ошибке.
}

// Цветовые коды ANSI для оформления вывода в консоль.
const (
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Purple = "\033[35m"
	Reset  = "\033[0m"
)

// Метки уровней логирования с использованием цветового оформления.
const (
	INFO    string = Green + "INFO" + Reset     // Информационные сообщения.
	WARNING        = Yellow + "WARNING" + Reset // Предупреждения.
	ERROR          = Red + "ERROR" + Reset      // Ошибки выполнения.
	FATAL          = Purple + "FATAL" + Reset   // Критические ошибки, приводящие к завершению программы.
)

// PrintInfo выводит информационное сообщение в консоль.
func PrintInfo(msg string) {
	log.Printf("[%s]: %s\n", INFO, msg)
}

// PrintWarning выводит предупреждение в консоль.
func PrintWarning(msg string) {
	log.Printf("[%s]: %s\n", WARNING, msg)
}

// PrintError выводит сообщение об ошибке в консоль.
func PrintError(msg string) {
	log.Printf("[%s]: %s\n", ERROR, msg)
}

// PrintFatal выводит критическое сообщение и завершает выполнение программы.
func PrintFatal(msg string) {
	log.Fatalf("[%s]: %s\n", FATAL, msg)
}

// SendError отправляет JSON-ответ с описанием ошибки клиенту.
// Устанавливает HTTP-заголовок Content-Type в "application/json".
func SendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
}
