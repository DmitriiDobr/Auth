package main

import (
	"auth/internal/auth/handler"
	repository "auth/internal/auth/repository"
	"auth/internal/auth/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	app := fiber.New()
	fmt.Println("Staring server...")
	if err := initConfig(); err != nil {
		fmt.Println("Ошибка при инициализации!")
	}
	config := repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Password: viper.GetString("db.password"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
	db, _ := config.InitDb()
	repo := repository.NewAuthService(db)
	authService := service.NewAuthService(repo)
	handlerAuth := handler.NewHandler(authService)
	handlerAuth.RegisterHandlers(app)
	err := app.Listen(viper.GetString("port"))
	if err != nil {
		return
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
