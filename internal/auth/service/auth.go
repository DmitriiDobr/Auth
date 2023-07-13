package service

import (
	"auth/internal/auth/types"
	"context"
	"crypto/sha1"
	"fmt"
)

// там где используется интерфейс там и используется.
// интерфейс репозитория
type Irepository interface {
	CreateUser(ctx context.Context, user, password string) (int, error)
	GetUser(ctx context.Context, username string) (*types.User, error)
}

type AuthService struct {
	repository Irepository
	salt       string
}

func NewAuthService(repository Irepository, salt string) *AuthService {
	return &AuthService{repository: repository, salt: salt}
}

func (a *AuthService) CreateUser(ctx context.Context, user types.User) (int, error) {
	password := GeneratePasswordHash(user.Password, a.salt)
	return a.repository.CreateUser(ctx, user.Username, password)
}

func (a *AuthService) GetUser(ctx context.Context, username string) (*types.User, error) {
	return a.repository.GetUser(ctx, username)
}

func GeneratePasswordHash(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
