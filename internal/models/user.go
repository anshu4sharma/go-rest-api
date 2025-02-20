package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	Password     string     `gorm:"type:varchar(255)" json:"-" validate:"required,password"`
	Name         string     `gorm:"type:varchar(100)" json:"name" validate:"required,name"`
	GoogleID     *string    `gorm:"type:varchar(100);unique;default:null" json:"google_id,omitempty"`
	Role         string     `gorm:"type:varchar(20);default:'user'" json:"role" validate:"required,oneof=user admin"`
	LastLoginAt  *time.Time `gorm:"default:null" json:"last_login_at,omitempty"`
	RefreshToken string     `gorm:"type:varchar(255)" json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Name     string `json:"name" validate:"required,name"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
