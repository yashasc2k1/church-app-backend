package utils

import (
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Set up a formatter that includes date, time, file, and line number
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,                  // Include timestamp
		TimestampFormat: "2006-01-02 15:04:05", // Custom timestamp format
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", filename + ":" + string(f.Line)
		},
	})

	Log.SetReportCaller(true)       // This enables the file name and line number
	Log.SetOutput(os.Stdout)        // Output to stdout
	Log.SetLevel(logrus.DebugLevel) // Set the logging level to debug for detailed logs
}
