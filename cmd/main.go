package main

import (
	"auth/internal/auth/handler"
	"auth/internal/auth/repository"
	"auth/internal/auth/service"
	"fmt"
	kafkaNotification "github.com/DmitriiDobr/kafkaNotification/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func main() {
	app := fiber.New()
	db, clientKafka, salt, jwtKey, err := run()
	repo := repository.NewAuthService(db)
	authService := service.NewAuthService(repo, salt)
	handlerAuth := handler.NewHandler(authService, clientKafka, salt, jwtKey)
	handlerAuth.RegisterHandlers(app)
	port := viper.GetString("port")
	err = app.Listen(port)
	if err != nil {
		panic(fmt.Sprintf("Server Listener does not work with port %s", port))
	}

}

func run() (*sqlx.DB, *kafkaNotification.Client, string, string, error) {
	if err := initConfig(); err != nil {
		panic("Ошибка при инициализации!")
	}
	configDB := repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Password: viper.GetString("db.password"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
	address := viper.GetString("kafka.address")
	configKafka := kafkaNotification.Config{
		Brokers: address,
		Topic:   viper.GetString("kafka.topic"),
		Address: address,
	}
	salt := viper.GetString("auth.salt")
	jwtKey := viper.GetString("auth.jwtKey")

	clientKafka, err := kafkaNotification.New(&configKafka)
	if err != nil {
		return nil, nil, "", "", err
	}
	db, err := configDB.InitDb()
	if err != nil {
		return nil, nil, "", "", err
	}
	return db, clientKafka, salt, jwtKey, nil

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
