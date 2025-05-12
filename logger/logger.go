package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Log is the shared global logger instance
var Log *logrus.Logger

// InitLogger initializes the logger based on the application environment
func InitLogger(env string) {
	Log = logrus.New()

	// Output ke stdout
	Log.SetOutput(os.Stdout)

	// Set log level (default: Info)
	Log.SetLevel(logrus.InfoLevel)

	// Use JSONFormatter in production, TextFormatter in development
	if strings.ToLower(env) == "production" {
		Log.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: false,
			TimestampFormat:  "2006-01-02T15:04:05Z07:00",
		})
	} else {
		// Development: readable format
		Log.SetLevel(logrus.DebugLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
			PadLevelText:  true,
		})
	}
}