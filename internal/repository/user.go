package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/newnorthblog/backend/internal/db"
	"github.com/newnorthblog/backend/internal/domain"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func newUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
	INSERT INTO "user"
	(id, username, email, "password")
	VALUES($1, $2, $3, $4);
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		if db.IsDuplicate(err) {
			return domain.ErrDuplicateEntry
		}
		return fmt.Errorf("insert user failed: %w", err)
	}

	return nil
}
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
	SELECT id, username, email, "password", created_at, updated_at, deleted_at
	FROM "user"
	WHERE email = $1;
	`

	var user domain.User
	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("select user failed: %w", err)
	}

	return &user, nil
}
