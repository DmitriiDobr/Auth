package service

import (
	repository "auth/internal/auth/repository"
	"context"
	"crypto/sha1"
	"fmt"
)

const salt = "nwueewui3212wlwo"

type AuthMethods interface {
	CreateUser(ctx context.Context, user repository.User) (int, error)
	GetUser(ctx context.Context, username string) (*repository.User, error)
}

type AuthService struct {
	repository *repository.AuthRepostitory
}

func NewAuthService(repository *repository.AuthRepostitory) *AuthService {
	return &AuthService{repository: repository}
}

func (a *AuthService) CreateUser(ctx context.Context, user repository.User) (int, error) {
	password := GeneratePasswordHash(user.Password)
	return a.repository.CreateUser(ctx, user.Username, password)
}

func (a *AuthService) GetUser(ctx context.Context, username string) (*repository.User, error) {
	return a.repository.GetUser(ctx, username)
}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
