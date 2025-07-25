package main

import (
	"net/http"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/routes"
)

func main() {
	config.SetEnv()
	http.HandleFunc("/", routes.HandleRoutes)
	http.ListenAndServe(":8080", nil)
}
