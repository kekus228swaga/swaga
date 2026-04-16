package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kekus228swaga/orderflow/auth-service/internal/domain/user"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash string) (*user.User, error) {
	var u user.User
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at`
	err := r.pool.QueryRow(ctx, query, email, passwordHash).Scan(&u.ID, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	err := r.pool.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil // user not found
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
