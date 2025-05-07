package auth

import (
	"context"
	"time"

	"simple-finance/internal/db"
	"simple-finance/internal/errs"
	"simple-finance/internal/tokens"
	"simple-finance/pkg/hash"
)

type Manager struct {
	db           *db.FinanceDB
	hasher       hash.PasswordHasher
	tokenManager *tokens.TokenManager
}

func NewManager(
	db *db.FinanceDB,
	hasher hash.PasswordHasher,
	tokenManager *tokens.TokenManager,
) *Manager {
	return &Manager{
		db:           db,
		hasher:       hasher,
		tokenManager: tokenManager,
	}
}

func (m *Manager) ComparePassword(ctx context.Context, userName, inputPass string) (string, error) {
	inputHashPass, err := m.hasher.Hash(inputPass)
	if err != nil {
		return "", err
	}

	userInfo, err := m.db.GetUserInfo(ctx, userName)
	if err != nil {
		return "", err
	}

	if userInfo.Password != inputHashPass {
		return "", errs.ErrInvalidPassword
	}

	return userInfo.ID, nil
}

func (m *Manager) MakeTokens(userID string, accessTokenTTL, refreshTokenTTL time.Duration) (string, string, error) {
	accessToken, err := m.tokenManager.NewJWT(tokens.TokenInfo{UserID: userID}, accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := m.tokenManager.NewJWT(tokens.TokenInfo{UserID: userID}, refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

func (m *Manager) RefreshTokens(refreshToken string, accessTokenTTL, refreshTokenTTL time.Duration) (string, string, error) {
	tokenInfo, err := m.tokenManager.Parse(refreshToken)
	if err != nil {
		return "", "", err
	}

	return m.MakeTokens(tokenInfo.UserID, accessTokenTTL, refreshTokenTTL)
}
