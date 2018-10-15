package main

import (
	"log"
	"net/http"

	"github.com/ericmcbride/go-dfw-testing/pkg/logging"
	server "github.com/ericmcbride/go-dfw-testing/pkg/server"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("logging", "DEBUG")
	viper.SetEnvPrefix("GO-DFW-TESTING")
	viper.AutomaticEnv()

	logging.ConfigureLogger(
		viper.Get("logging").(string),
	)

	handler := server.New()
	log.Println("Starting server on: ", ":8080")

	log.Fatal(http.ListenAndServe(":8080", handler))
}
