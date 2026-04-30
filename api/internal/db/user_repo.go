package db

import (
	"context"
	"database/sql"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, u models.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, name, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, u.ID, u.Email, u.PasswordHash, u.Name, u.CreatedAt)
	return err
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, name, created_at FROM users WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt)
	return u, err
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, name, created_at FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt)
	return u, err
}
