package models

import "time"

// Tokens represents authentication tokens
// @Description  Authentication tokens response
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SignInInput represents user login credentials
// @Description  User login credentials
type SignInInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// SignUpInput represents user registration data
// @Description  User registration data
type SignUpInput struct {
	Email    string `json:"email" validate:"required"`
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserInfo represents user information
// @Description  User information
type UserInfo struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	UserName  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// RefreshInput represents refresh token request
// @Description  Refresh token request
type RefreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
