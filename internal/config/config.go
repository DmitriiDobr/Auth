package config

import "github.com/spf13/viper"

type Config struct {
	Database struct {
		Host     string
		Port     string
		Password string
		Username string
		DBName   string
		SSLMode  string
		Driver   string
	}
	Kafka struct {
		Address string
		Topic   string
	}
	Auth struct {
		JwtKey string
		Salt   string
	}
}

func NewConfig() *Config {
	if err := initConfig(); err != nil {
		panic("Ошибка при инициализации!")
	}
	cfg := &Config{
		Database: struct {
			Host     string
			Port     string
			Password string
			Username string
			DBName   string
			SSLMode  string
			Driver   string
		}{Host: viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Password: viper.GetString("db.password"),
			Username: viper.GetString("db.username"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
			Driver:   viper.GetString("db.driver"),
		},
		Kafka: struct {
			Address string
			Topic   string
		}{Address: viper.GetString("kafka.address"),
			Topic: viper.GetString("kafka.topic")},
		Auth: struct {
			JwtKey string
			Salt   string
		}{JwtKey: viper.GetString("auth.jwtKey"), Salt: viper.GetString("auth.salt")},
	}
	return cfg
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
