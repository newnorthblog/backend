package repository

import (
	"context"

	"github.com/newnorthblog/backend/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Users
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users: newUserRepository(db),
	}
}

type Users interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
