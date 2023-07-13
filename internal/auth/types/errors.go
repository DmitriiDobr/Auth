package types

import "errors"

var ErrUnauthorized = errors.New("Пользователь не авторизован!")

var ErrTokenCreation = errors.New("Не получилось создать токен!")

var ErrTimeRunnedOut = errors.New("Время токена просрочено")

var ErrUserCreation = errors.New("Пользователя не удалось создать")

var ErrCantSendMessage = errors.New("Сообщение не удалось отправить пользователю")

var ErrIncorrectPassword = errors.New("Ваш пароль некорректен!")
