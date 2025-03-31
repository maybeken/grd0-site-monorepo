package utils

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func InitiateLogger() *logrus.Logger {
	var log_level logrus.Level

	switch GetEnvWithFallback("LOG_LEVEL", "WARN") {
	case "DEBUG":
		log_level = logrus.DebugLevel
	case "INFO":
		log_level = logrus.InfoLevel
	case "WARN":
		log_level = logrus.WarnLevel
	case "ERROR":
		log_level = logrus.WarnLevel
	default:
		log_level = logrus.DebugLevel
	}

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log_level)

	return log
}

func EchoLogSkipper(c echo.Context) bool {
	// Skip health check endpoint
	return c.Request().URL.Path == "/health"
}

func EchoLogValueFunc(c echo.Context, v middleware.RequestLoggerValues) error {
	log := InitiateLogger()

	log.WithFields(logrus.Fields{
		"Status":    v.Status,
		"RealIP":    c.RealIP(),
		"URI":       v.URI,
		"RequestID": c.Response().Header().Get(echo.HeaderXRequestID),
	}).Infof("HTTP Request[%s]", v.Latency)

	return nil
}
