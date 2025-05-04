package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenInfo struct {
	UserID string
}

type TokenManager struct {
	signingKey string
}

func NewTokenManager(signingKey string) (*TokenManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &TokenManager{signingKey: signingKey}, nil
}

func (m *TokenManager) NewJWT(tokenInfo TokenInfo, ttl time.Duration) (string, error) {
	if tokenInfo.UserID == "" {
		return "", fmt.Errorf("userID or user role is empty")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(ttl).Unix(),
		"sub": tokenInfo.UserID,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *TokenManager) Parse(accessToken string) (TokenInfo, error) {
	keyFunc := func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	}
	token, err := jwt.Parse(accessToken, keyFunc)
	if err != nil {
		return TokenInfo{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenInfo{}, fmt.Errorf("error get user claims from token")
	}

	return TokenInfo{
		UserID: claims["sub"].(string),
	}, nil
}
