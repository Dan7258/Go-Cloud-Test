package models

type RateLimits struct {
	ClientID   string `json:"client_id" gorm:"primary_key;"`
	Capacity   int    `json:"capacity" gorm:"not null;"`
	RatePerSec int    `json:"rate_per_sec" gorm:"not null;"`
}

func CreateClient(rl RateLimits) error {
	err := DB.Create(&rl).Error
	return err
}

func UpdateClient(rl *RateLimits) error {
	err := DB.Model(new(RateLimits)).Where("client_id = ?", rl.ClientID).Updates(rl).Error
	return err
}

func DeleteClient(clientID string) error {
	err := DB.Where("client_id = ?", clientID).Delete(new(RateLimits)).Error
	return err
}

func GetClient(clientID string) (RateLimits, error) {
	client := new(RateLimits)
	err := DB.Where("client_id = ?", clientID).First(client).Error
	return *client, err
}

func ThsClientExists(clientID string) bool {
	client := new(RateLimits)
	result := DB.Select("client_id").Where("client_id = ?", clientID).First(client)
	if result.Error != nil {
		return false
	}
	return true
}
