package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id       uint64 `json:"id" gorm:"primaryKey"`
	Email    string `gorm:"unique,not null"`
	Password string `gorm:"not null"`
}
