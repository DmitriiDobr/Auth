package main

import (
	"auth/internal/auth/handler"
	"auth/internal/auth/hashword"
	"auth/internal/auth/notification"
	"auth/internal/auth/repository"
	"auth/internal/auth/service"
	"auth/internal/config"
	"auth/pkg/logging"
	"auth/pkg/postgresLogger"
	"fmt"
	kafkaNotification "github.com/DmitriiDobr/kafkaNotification/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

func main() {
	app := fiber.New()
	cfg := config.NewConfig()
	logging.InitLogger()
	db, err := InitDb(cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.Password,
		cfg.Database.Username,
		cfg.Database.SSLMode,
		cfg.Database.Driver)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Не получилось инициализировать базу данных %s", err.Error()))
		panic(err)
	}
	kafka, err := InitKafkaBroker(cfg.Kafka.Address,
		cfg.Kafka.Topic,
		cfg.Kafka.Address)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Не получилось инициализировать кафку %s", err.Error()))
	}
	generatePass := hashword.NewHash(cfg.Auth.Salt)
	dbWithLog := postgresLogger.NewDbWithLogger(db, logrus.New())
	repo := repository.NewAuthRepository(dbWithLog)
	authService := service.NewAuthService(repo, kafka, generatePass, cfg.Auth.JwtKey)
	handlerAuth := handler.NewHandler(authService)
	handlerAuth.RegisterHandlers(app)
	port := viper.GetString("port")
	err = app.Listen(port)
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Server Listener does not work with port %s %s", port,
			err.Error()))
		panic(fmt.Sprintf("Server Listener does not work with port %s", port))
	}

}

func InitDb(host, port, dbName, password, username, sslmode, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, username, dbName, password, sslmode))
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Не удалось подключится к базе!")
		return nil, err
	}
	return db, db.Ping()
}

func InitKafkaBroker(brokers, topic, address string) (*notification.Notification, error) {
	cfg := &kafkaNotification.Config{
		Brokers: brokers,
		Topic:   topic,
		Address: address,
	}
	kafka, _ := kafkaNotification.New(cfg)
	return notification.NewNotification(kafka), nil
}
