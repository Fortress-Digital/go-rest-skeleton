package model

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	Email    string `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
}
