package main

import (
	"net/http"

	"github.com/gocroot/core/routes"
)

func main() {
	http.HandleFunc("/", routes.HandleRoutes)
	http.ListenAndServe(":8080", nil)
}
