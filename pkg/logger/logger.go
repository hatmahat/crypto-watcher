package logger

import (
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func MiddlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		makeLogEntry(c, start, err, "Incoming Request")
		return err
	}
}

func LogInfoWithCustomTime(message string) {
	logrus.WithFields(logrus.Fields{
		"at": time.Now().Format("2006-01-02 15:04:05"),
	}).Info(message)
}

func LogDebugWithCustomTime(message string) {
	logrus.WithFields(logrus.Fields{
		"at": time.Now().Format("2006-01-02 15:04:05"),
	}).Debug(message)
}

func makeLogEntry(c echo.Context, start time.Time, err error, message string) {
	fields := logrus.Fields{
		"method":  c.Request().Method,
		"uri":     c.Request().URL.String(),
		"ip":      c.Request().RemoteAddr,
		"latency": time.Since(start).String(),
	}

	if err != nil {
		logrus.WithFields(fields).WithField("error", err.Error()).Warn(message)
		return
	}
	logrus.WithFields(fields).Info(message)
}
