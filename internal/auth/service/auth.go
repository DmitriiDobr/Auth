package service

import (
	"auth/internal/auth/types"
	errors2 "auth/pkg/errors"
	"context"
	errUsual "errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// там где используется интерфейс там и используется.
// интерфейс репозитория
type UserAction interface {
	CreateUser(ctx context.Context, user, password string) (int, error)
	GetUser(ctx context.Context, username string) (*types.User, error)
}

type Broker interface {
	Notify(ctx context.Context, message types.Message) error
}

type CreatePassword interface {
	Generate(password string) string
}

type AuthService struct {
	repository       UserAction
	messageBroker    Broker
	passwordCreation CreatePassword
	jwtKey           string
}

func NewAuthService(repository UserAction, messageBroker Broker, passwordCreator CreatePassword, jwtKey string) *AuthService {
	return &AuthService{repository: repository,
		messageBroker:    messageBroker,
		passwordCreation: passwordCreator,
		jwtKey:           jwtKey}
}

func (a *AuthService) Login(ctx context.Context, creds types.Login, header, body string) (time.Time, *jwt.Token, string, error) {
	user, _ := a.repository.GetUser(ctx, creds.Username)
	if user.Password != a.passwordCreation.Generate(creds.Password) {
		return time.Now(), nil, "", errors2.Cause().BadRequest().New(
			"Пароль не соответсвует паролю пользователя!")
	}
	expirationDate := time.Now().Add(30 * time.Second)
	claims := &types.Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.jwtKey))
	if err != nil {
		return time.Now(), nil, "", errors2.Cause().InternalServerError().New("Не получилось создать токен!")
	}
	msg := types.Message{
		UserID: user.Id,
		Status: types.Success,
		Header: header,
		Body:   body,
	}

	err = a.messageBroker.Notify(ctx, msg)
	if err != nil {
		return time.Now(), nil, "", errors2.Cause().InternalServerError().
			New("Сообщение не удалось отправить пользователю!")
	}
	return expirationDate, token, tokenString, err

}

func (a *AuthService) Register(ctx context.Context, user types.User, header, body string) error {
	password := a.passwordCreation.Generate(user.Password)
	_, err := a.repository.CreateUser(ctx, user.Username, password)
	if err != nil {
		return errors2.Cause().InternalServerError().
			New("Пользователя не удалось создать!")
	}

	msg := types.Message{
		UserID: user.Id,
		Status: types.Success,
		Header: header,
		Body:   body,
	}

	err = a.messageBroker.Notify(ctx, msg)
	if err != nil {
		return errors2.Cause().InternalServerError().
			New("Сообщение не удалось отправить пользователю!")
	}
	return nil
}

func (a *AuthService) Refresh(token string) (time.Time, string, error) {
	if token == "mistake" {
		return time.Now(), "", errors2.Cause().Unauthorized().
			New("Пользователь не авторизован в системе!")
	}
	claims := &types.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})

	if err != nil {
		if errUsual.Is(err, jwt.ErrSignatureInvalid) {
			return time.Now(), "", errors2.Cause().Unauthorized().
				New("Пользователь не авторизован в системе!")
		}
		return time.Now(), "", err
	}
	if !tkn.Valid {
		return time.Now(), "", errors2.Cause().Unauthorized().
			New("Пользователь не авторизован в системе!")
	}
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return time.Now(), "", errors2.Cause().Unauthorized().
			New("Время токена просрочено!")

	}
	//конфиг
	expirationalTime := time.Now().Add(30 * time.Second)
	claims.ExpiresAt = jwt.NewNumericDate(expirationalTime)
	tokenUpdate := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenUpdate.SignedString(a.jwtKey)
	if err != nil {
		return time.Now(), "", errors2.Cause().InternalServerError().
			New("Не получилось создать токен!")
	}
	return expirationalTime, tokenString, err

}
