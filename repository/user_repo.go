package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"travel-planning/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := models.User{}

	query := `SELECT user_id, first_name, last_name, email, password_hash, created_at
	          FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Insert(user *models.User, password string) (int, error) {
	existingUser, err := r.GetByEmail(user.Email)
	if err != nil {
		return 0, fmt.Errorf("db error during existence check:%w", err)
	}

	if existingUser != nil {
		return existingUser.UserID, fmt.Errorf("user with email %s already exists", user.Email)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashedPass)

	query := `INSERT INTO users (first_name, last_name, email, password_hash, created_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING user_id`

	var userID int
	currTime := time.Now()
	err = r.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		currTime,
	).Scan(&userID)

	if err != nil {
		log.Printf("SQL Error inserting user: %v", err)
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return userID, nil
}
