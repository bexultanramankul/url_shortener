package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func InitLogger() {
	Log.SetFormatter(&logrus.TextFormatter{
		DisableColors:             false,                 // Отключить цвета в логах
		ForceColors:               true,                  // Принудительное включение цветов
		DisableTimestamp:          false,                 // Отключить вывод времени
		FullTimestamp:             true,                  // Полное время вместо относительного
		TimestampFormat:           "2006-01-02 15:04:05", // Формат времени (YYYY-MM-DD HH:MM:SS)
		DisableSorting:            false,                 // Отключить сортировку полей (по умолчанию false)
		DisableQuote:              true,                  // Убрать кавычки у строк в логе
		ForceQuote:                false,                 // Принудительно добавлять кавычки к строкам
		EnvironmentOverrideColors: true,                  // Использовать цвета из переменных окружения
		QuoteEmptyFields:          false,                 // Заключать пустые поля в кавычки
	})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
}
