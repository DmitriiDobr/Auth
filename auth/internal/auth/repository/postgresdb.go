package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Config struct {
	Host     string
	Port     string
	Password string
	Username string
	DBName   string
	SSLMode  string
}

func (c *Config) InitDb() (*sqlx.DB, error) {
	//db, err := sqlx.Open("postgres", "postgres://auth:qwerty@localhost:5436/postgres?sslmode=disable")
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.DBName, c.Password, c.SSLMode))
	if err != nil {
		log.Fatal("Не удалось подключится к базе!")
		return nil, err
	}
	return db, db.Ping()
}
