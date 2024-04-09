package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vekshinnikita/golang_music"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db}
}

func (r *AuthPostgres) CreateUser(user golang_music.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username string, password string) (golang_music.User, error) {
	var user golang_music.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 and password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
