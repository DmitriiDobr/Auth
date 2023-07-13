package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitLogger() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.WarnLevel)
	logger.SetReportCaller(true)
}

func LoginLogger(username, header, body string) *logrus.Entry {
	loginLogger := logrus.WithFields(logrus.Fields{"username": username,
		"header": header, "body": body})
	return loginLogger
}

func RefreshLogger(token string) *logrus.Entry {
	refreshLogger := logrus.WithFields(logrus.Fields{"token": token, "action": "Обновление токена"})
	return refreshLogger
}

func RegisterLogger(username, password string) *logrus.Entry {
	registerLogger := logrus.WithFields(logrus.Fields{"username": username,
		"password": password})
	return registerLogger

}

func GetUserLogger(username string) *logrus.Entry {
	getUserLogger := logrus.WithFields(logrus.Fields{"username": username})
	return getUserLogger
}

func CreateUserLogger(username, password string) *logrus.Entry {
	createUserLogger := logrus.WithFields(logrus.Fields{"username": username,
		"password": password})
	return createUserLogger
}
