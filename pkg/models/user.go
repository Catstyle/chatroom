package models

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex"`
	TokenHash string
	CTime     int    `gorm:"autoCreateTime"`
}
