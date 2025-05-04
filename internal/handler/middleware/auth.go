package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"simple-finance/internal/handler/response"
	"simple-finance/internal/tokens"
)

type ContextKey int

const (
	authorizationHeader string     = "Authorization"
	TokenInfoKey        ContextKey = 1
)

type AuthMiddleware struct {
	tokenManager *tokens.TokenManager
}

func NewAuthMiddleware(tokenManager *tokens.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
	}
}

func (h *AuthMiddleware) MakeAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authorizationHeader)
		tokenInfo, err := h.parseAuthHeader(authHeader)

		if err != nil {
			response.Unauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), TokenInfoKey, tokenInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *AuthMiddleware) parseAuthHeader(authHeader string) (tokens.TokenInfo, error) {
	if authHeader == "" {
		return tokens.TokenInfo{}, errors.New("empty auth header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return tokens.TokenInfo{}, errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return tokens.TokenInfo{}, errors.New("token is empty")
	}
	accessToken := headerParts[1]

	return h.tokenManager.Parse(accessToken)
}
