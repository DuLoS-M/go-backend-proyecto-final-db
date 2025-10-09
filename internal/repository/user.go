package repository

import (
	"database/sql"
	"fmt"

	"proyecto-bd-final/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}
	return nil
}
