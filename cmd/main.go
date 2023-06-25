package main

import (
	"auth/internal/auth/handler"
	repository "auth/internal/auth/repository"
	"auth/internal/auth/service"
	"fmt"
	kafkaNotification "github.com/DmitriiDobr/kafkaNotification/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func main() {
	app := fiber.New()
	fmt.Println("Staring server...")
	if err := initConfig(); err != nil {
		fmt.Println("Ошибка при инициализации!")
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

	clientKafka, err := kafkaNotification.New(&configKafka)
	if err != nil {
		fmt.Println(err)
		return
	}
	db, _ := configDB.InitDb()
	repo := repository.NewAuthService(db)
	authService := service.NewAuthService(repo)
	handlerAuth := handler.NewHandler(authService, clientKafka)
	handlerAuth.RegisterHandlers(app)
	err = app.Listen(viper.GetString("port"))
	if err != nil {
		fmt.Println("nanan")
		return
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
