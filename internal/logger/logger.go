package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Level string `envconfig:"APP_LOG_LEVEL" default:"debug"`
}

func Setup(cfg Config) (err error) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(cfg.Level)
	if err != nil {
		return
	}
	log.SetLevel(level)
	return
}
