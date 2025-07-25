package main

import (
	"net/http"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/routes"
)

func main() {
	config.SetEnv()
	http.HandleFunc("/", routes.HandleRoutes)
	err := http.ListenAndServe(":3000", nil)
	println(err.Error())
}
