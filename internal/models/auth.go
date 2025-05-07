package models

import "time"

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignInInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpInput struct {
	Email    string `json:"email" validate:"required"`
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserInfo struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	UserName  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
