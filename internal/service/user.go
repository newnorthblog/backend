package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/newnorthblog/backend/internal/domain"
	"github.com/newnorthblog/backend/internal/pkg/tokenmanager"
	"github.com/newnorthblog/backend/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository repository.Users
	logger         *slog.Logger
	tokenManager   *tokenmanager.Manager
}

func newUserService(
	userRepository repository.Users,
	logger *slog.Logger,
	tokenManager *tokenmanager.Manager,
) *userService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
		tokenManager:   tokenManager,
	}
}

type RegisterInput struct {
	Email    string
	Username string
	Password string
}

func (s *userService) Register(ctx context.Context, input *RegisterInput) error {
	userID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate user id failed: %w", err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword failed: %w", err)
	}

	if err := s.userRepository.Create(ctx, &domain.User{
		ID:       userID,
		Username: input.Username,
		Email:    input.Email,
		Password: passHash,
	}); err != nil {
		if errors.Is(err, domain.ErrDuplicateEntry) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("create user failed: %w", err)
	}

	return nil
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *userService) Login(ctx context.Context, email, password string) (*Tokens, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email failed: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return nil, ErrUserInvalidCredentials
	}

	accessToken, _, err := s.tokenManager.NewJWT(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate access token failed: %w", err)
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: "",
	}, nil
}
