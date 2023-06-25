package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthDbManipulation interface {
	CreateUser(ctx context.Context, username, password string) (int, error)
	GetUser(ctx context.Context, username string) (*User, error)
}

type AuthRepostitory struct {
	db *sqlx.DB
	AuthDbManipulation
}

func NewAuthService(db *sqlx.DB) *AuthRepostitory {
	return &AuthRepostitory{db: db}
}

func (a *AuthRepostitory) GetUser(ctx context.Context, username string) (*User, error) {
	var userData User
	query := fmt.Sprintf(`
SELECT u.id,u.username,u.password
    FROM user_auth u 
WHERE username=$1 
`)
	if err := a.db.GetContext(ctx, &userData, query, username); err != nil {
		return nil, err
	}
	return &userData, nil

}

func (a *AuthRepostitory) CreateUser(ctx context.Context, username, password string) (idUser int, err error) {
	userRegistration := Registrator{username, password}
	query := fmt.Sprintf(`
INSERT into 
    user_auth (username, password)
values (:username,:password) RETURNING id`)
	query, args, err := sqlx.Named(query, userRegistration)
	query = a.db.Rebind(query)
	return idUser, a.db.GetContext(ctx, &idUser, query, args...)
}
