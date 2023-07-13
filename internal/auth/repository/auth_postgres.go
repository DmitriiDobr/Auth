package repository

import (
	"auth/internal/auth/types"
	"context"
	"fmt"
)

type DbLogger interface {
	GetUserWithLogger(ctx context.Context, user interface{}, query string, args ...interface{}) error
	CreateUserWithLogger(registrator types.Registrator, query string) (string, []interface{}, error)
}

// Логирование запроса к БД
type AuthRepostitory struct {
	db DbLogger
}

func NewAuthRepository(db DbLogger) *AuthRepostitory {
	return &AuthRepostitory{db: db}
}

func (a *AuthRepostitory) GetUser(ctx context.Context, username string) (*types.User, error) {
	var userData types.User
	query := fmt.Sprintf(`
SELECT u.id,u.username,u.password
    FROM user_auth u 
WHERE username=$1 
`)
	err := a.db.GetUserWithLogger(ctx, &userData, query, username)
	return &userData, err

}

// библиотека обертка над sqlx, вызов логрус
func (a *AuthRepostitory) CreateUser(ctx context.Context, username, password string) (idUser int, err error) {
	userRegistration := types.Registrator{username, password}
	query := fmt.Sprintf(`
INSERT into 
    user_auth (username, password)
values (:username,:password) RETURNING id`)
	//query, args, err := a.db.CreateUserWithLogger(userRegistration, query)
	return idUser, a.db.GetUserWithLogger(ctx, &idUser, query, args...)

}
