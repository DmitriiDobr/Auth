package postgresLogger

import (
	"auth/internal/auth/types"
	"context"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type DbWithLogger struct {
	db  *sqlx.DB
	log *logrus.Logger
}

func NewDbWithLogger(db *sqlx.DB, log *logrus.Logger) *DbWithLogger {
	return &DbWithLogger{db: db, log: log}
}

// GET CONTEXT WITH LOGGER context dest, query, arg (как any)
func (pg *DbWithLogger) GetUserWithLogger(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := pg.db.GetContext(ctx, &dest, query, args)
	if err != nil {
		pg.log.Errorf("Пользователь c именем %s не найден", args)
		return err
	}
	pg.log.Infof("Пользователь c именем %s найден", args)
	return err

}

//GetContextWithLogger (значение вернутся)
// exec если ожидаем, что ничего не вернется

func (pg *DbWithLogger) CreateUserWithLogger(registrator types.Registrator, query string) (string, []interface{}, error) {
	begin := time.Now()

	if err := cache.New(database.GetDB()).TableTab().FullDelete(ctx); err != nil {
		return err
	}

	query, args, err := sqlx.Named(query, registrator)
	if err != nil {
		pg.log.Errorf("Пользователь c именем %s не создан. Ваш запрос - %s", registrator.Username, query)
		return "", args, err
	}

	query = pg.db.Rebind(query)
	_, err = pg.db.Exec(query, args)
	if err != nil {
		return "", nil, err
	}
	pg.log.Infof(
		"Завершен процесс очистки кеша для позиций документов. Затраченное время: %s",
		time.Since(begin),
	)
	pg.log.Infof("Время - %s ,Пользователь c именем %s найден, Ваш запрос - %s", delta,
		registrator.Username,
		query)

	return query, args, err

}
