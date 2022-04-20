package models

import (
	"time"
)

type Token struct {
	Token     string    `json:"token" gorm:"primary_key"`
	Active    int16     `json:"active"` // To manually change status
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
