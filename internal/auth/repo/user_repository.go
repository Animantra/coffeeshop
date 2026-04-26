package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	model "github.com/thangchung/go-coffeeshop/internal/auth/domain"
)

// UserRepository handles all database operations for users.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Migrate creates the users table if it does not exist.
func (r *UserRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS auth_users (
		id            SERIAL PRIMARY KEY,
		username      VARCHAR(50)  NOT NULL UNIQUE,
		email         VARCHAR(255) NOT NULL UNIQUE,
		password_hash TEXT         NOT NULL,
		created_at    TIMESTAMP    NOT NULL DEFAULT NOW()
	);`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("migrate auth_users: %w", err)
	}
	return nil
}

// Create inserts a new user and returns it with the generated ID.
func (r *UserRepository) Create(username, email, passwordHash string) (*model.User, error) {
	user := &model.User{}
	query := `
		INSERT INTO auth_users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, created_at`

	err := r.db.QueryRow(query, username, email, passwordHash).
		Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

// FindByEmail looks up a user by email, including the password hash for verification.
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, username, email, password_hash, created_at
		FROM auth_users
		WHERE email = $1`

	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // caller checks for nil
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return user, nil
}
