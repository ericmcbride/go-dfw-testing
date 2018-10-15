package logging

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Fields = logrus.Fields

var Logger *logrus.Logger

// JSON-API error object
type JsonError struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Title   string `json:"title"`
}

// Configure logger. Call once at app initialization.
func ConfigureLogger(level string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if level == "DEBUG" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
	}
	Logger = logrus.StandardLogger()
}

// Return context logger populated with request information
func GetLog(ctx context.Context) *logrus.Entry {
	if Logger == nil {
		log.Fatal("Logger not enabled. Call Configure() first.")
	}

	if r := ctx.Value("request"); r != nil {
		return Logger.WithField("request", r)
	}
	return Logger.WithField("request", "")
}

// Send back nicely formatted JSON Error to Client thats sending requests
func FormatError(ctx context.Context, w http.ResponseWriter, status int, data JsonError) {
	log := GetLog(ctx)
	data.Status = strconv.Itoa(status)

	// If a message comes in nil, we default the message to the title
	if data.Message == "" {
		data.Message = data.Title
	}

	output, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("Unable to marshal error: ", err)
		return
	}

	http.Error(w, string(output), status)
}
