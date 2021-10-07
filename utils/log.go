package utils

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

// InitLogger initializes a global logger,
// also returns a writer for other callers to use (e.g. middleware)
func InitLogger() *io.PipeWriter {
	Log.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}) // Log as JSON instead of the default ASCII formatter.
	Log.SetOutput(os.Stdout)                                                        // Output to stdout instead of the default stderr
	Log.SetLevel(logrus.InfoLevel)                                                  // Only Log the info severity or above.
	Log.SetReportCaller(true)                                                       // Trace callers
	return Log.Writer()
}
