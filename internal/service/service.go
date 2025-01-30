package service

import (
	"context"
	"log/slog"

	"github.com/newnorthblog/backend/internal/config"
	"github.com/newnorthblog/backend/internal/pkg/tokenmanager"
	"github.com/newnorthblog/backend/internal/repository"
)

type Services struct {
	Users
}

type Deps struct {
	Logger       *slog.Logger
	Config       *config.Config
	Repos        *repository.Repositories
	TokenManager *tokenmanager.Manager
}

func NewServices(deps Deps) *Services {
	return &Services{
		Users: newUserService(deps.Repos.Users, deps.Logger, deps.TokenManager),
	}
}

type Users interface {
	Register(ctx context.Context, input *RegisterInput) error
	Login(ctx context.Context, email, password string) (*Tokens, error)
}
