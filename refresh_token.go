package gotwitter

import (
	"context"
	"time"
)

var (
	AccessTokenLifeTime  = time.Minute * 15   // 15 minutes
	RefreshTokenLifeTime = time.Hour * 24 * 7 // 7 days
)

type RefreshToken struct {
	ID        string
	Name      string
	UserID    string
	LastUseAt time.Time
	ExpiresAt time.Time
	CreatedAt time.Time
}

type CreateRefreshTokenParams struct {
	Sub  string
	Name string
}

type RefreshTokenRepo interface{ 
	Create(ctx context.Context, params CreateRefreshTokenParams) (RefreshToken, error)
	GetByID(ctx context.Context, id string) (RefreshToken, error)
}