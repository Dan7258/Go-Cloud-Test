package models

type RateLimits struct {
	ClientID      string `json:"client_id" gorm:"primary_key;"`
	Capacity      int64  `json:"capacity" gorm:"not null;"`
	RatePerSecond int64  `json:"rate_per_second" gorm:"not null;"`
}

func CreateClient(rl *RateLimits) error {
	err := DB.Create(&rl).Error
	return err
}

func UpdateClient(rl *RateLimits) error {
	err := DB.Model(new(RateLimits)).Where("client_id = ?", rl.ClientID).Updates(rl).Error
	return err
}

func DeleteClient(ID int64) error {
	err := DB.Delete(new(RateLimits), ID).Error
	return err
}

func GetClient(ID int64) (*RateLimits, error) {
	client := new(RateLimits)
	err := DB.First(client, ID).Error
	return client, err
}
