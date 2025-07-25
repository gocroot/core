package main

import (
	"net/http"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/routes"
)

func main() {
	config.SetEnv()
	println("Cek Environtment Variabel")
	println("Port : " + config.Port)
	println("Private Key : " + config.PrivateKey)
	println("Public Key : " + config.PublicKey)
	println("Mongostring : " + config.MongoString)

	http.HandleFunc("/", routes.HandleRoutes)
	err := http.ListenAndServe(config.Port, nil)
	println(err.Error())
}
