package service

import (
	"auth/internal/auth/mock"
	"auth/internal/auth/types"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/tjarratt/babble"
	"testing"
)

func TestRegister(t *testing.T) {
	//username, password := GenerateUserNameAndPassword()
	username := "dima"
	password := "bop"
	ctx := context.Background()
	user := types.User{
		Id:       1,
		Username: username,
		Password: password,
	}
	jwtKey := "my_secret_key"
	header := "user creation"
	body := "username password"

	hashed_password := "6e7775656577756933323132776c776fa2bcc90502334df57e081c9e7d5ac2d33cf31095"

	msg := types.Message{
		UserID: 1,
		Status: types.Success,
		Header: header,
		Body:   body,
	}

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	// структура init, которая реализует интерфейс party.NamesLister

	mockedUserAction := mock.NewMockUserAction(ctrl)
	kafkaBroker := mock.NewMockBroker(ctrl)
	mockedGenerator := mock.NewMockCreatePassword(ctrl)

	mockedUserAction.EXPECT().CreateUser(gomock.Any(), username, hashed_password).Return(1, nil)
	kafkaBroker.EXPECT().Notify(gomock.Any(), msg).Return(nil)
	mockedGenerator.EXPECT().Generate(password).Return(hashed_password)

	serve := NewAuthService(mockedUserAction, kafkaBroker, mockedGenerator, jwtKey)

	err := serve.Register(ctx, user, header, body)
	if err != nil {
		return
	}

	//service.GeneratePasswordHash()

}

func TestLogin(t *testing.T) {
	username := "dima"
	password := "bop"
	hashed_password := "6e7775656577756933323132776c776fa2bcc90502334df57e081c9e7d5ac2d33cf31095"
	header := "user login"
	body := "username password"
	ctx := context.Background()

	user := types.User{Username: username, Password: hashed_password}

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockedUserAction := mock.NewMockUserAction(ctrl)
	kafkaBroker := mock.NewMockBroker(ctrl)
	mockedGenerator := mock.NewMockCreatePassword(ctrl)

	mockedUserAction.EXPECT().GetUser(ctx, username).Return(user, nil)

}

func TestRefresh(t *testing.T) {

}

func GenerateUserNameAndPassword() (username, password string) {
	babbler := babble.NewBabbler()
	return babbler.Babble(), babbler.Babble()
}
