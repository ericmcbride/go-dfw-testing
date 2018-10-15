package server

import (
	"github.com/ericmcbride/go-dfw-testing/pkg/handlers"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Handler struct {
	*mux.Router
}

func New() http.Handler {
	m := mux.NewRouter()

	m.HandleFunc("/{health:health(?:\\/)?}", HealthEndpointHandler)
	m.HandleFunc("/{cars:cars(?:\\/)?}", handlers.CarsHandler)

	return m
}

func HealthEndpointHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"status": "OK"}`)
}
