package model

type Test struct {
	ID   uint   `json:"id" gorm:"primary_key;default:auto_random()"`
	Name string `json:"name" gorm:"type:varchar(255);not null"`
}
