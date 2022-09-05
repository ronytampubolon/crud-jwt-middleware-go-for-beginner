package models

import "time"

type User struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"created_at"`
}
