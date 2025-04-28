// Package models содержит логику для работы с клиентами и их лимитами в базе данных.
package models

// RateLimits представляет данные лимитов для клиента, такие как максимальная емкость и количество запросов в секунду.
type RateLimits struct {
	ClientID   string `json:"client_id" gorm:"primary_key;"` // Идентификатор клиента (первичный ключ).
	Capacity   int    `json:"capacity" gorm:"not null;"`     // Максимальная емкость клиента (например, максимальное количество запросов).
	RatePerSec int    `json:"rate_per_sec" gorm:"not null;"` // Максимальное количество запросов в секунду для клиента.
}

// CreateClient добавляет нового клиента в таблицу rate_limits.
// Возвращает ошибку, если создание записи не удалось.
func CreateClient(rl RateLimits) error {
	err := DB.Create(&rl).Error
	return err
}

// UpdateClient обновляет данные клиента в таблице rate_limits по его идентификатору.
// Принимает указатель на объект RateLimits для изменения существующей записи.
func UpdateClient(rl *RateLimits) error {
	err := DB.Model(new(RateLimits)).Where("client_id = ?", rl.ClientID).Updates(rl).Error
	return err
}

// DeleteClient удаляет клиента из таблицы rate_limits по его идентификатору.
// Возвращает ошибку, если удаление не удалось.
func DeleteClient(clientID string) error {
	err := DB.Where("client_id = ?", clientID).Delete(new(RateLimits)).Error
	return err
}

// GetClient извлекает данные клиента из таблицы rate_limits по его идентификатору.
// Возвращает структуру RateLimits и ошибку, если клиент не найден.
func GetClient(clientID string) (RateLimits, error) {
	client := new(RateLimits)
	err := DB.Where("client_id = ?", clientID).First(client).Error
	return *client, err
}

// ThsClientExists проверяет, существует ли клиент с данным идентификатором в базе данных.
// Возвращает true, если клиент найден, и false, если нет.
func ThsClientExists(clientID string) bool {
	client := new(RateLimits)
	result := DB.Select("client_id").Where("client_id = ?", clientID).First(client)
	if result.Error != nil {
		return false
	}
	return true
}
