package logger

import (
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func MiddlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		makeLogEntry(c, start, err).Info("incoming request")
		return next(c)
	}
}

func LogWithCustomTime(message string) {
	log.WithFields(log.Fields{
		"at": time.Now().Format("2006-01-02 15:04:05"),
	}).Info(message)
}

func makeLogEntry(c echo.Context, start time.Time, err error) *log.Entry {
	entry := log.WithFields(log.Fields{
		"method":  c.Request().Method,
		"uri":     c.Request().URL.String(),
		"ip":      c.Request().RemoteAddr,
		"latency": time.Since(start).String(),
	})

	if err != nil {
		entry = entry.WithField("error", err.Error())
	}
	return entry
}
