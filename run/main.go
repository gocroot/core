package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/routes"
)

func main() {
	config.SetEnv()
	log.Println("Cek Environment Variabel")
	log.Println("PORT        :", config.Port)
	log.Println("PRIVATEKEY  :", status(config.PrivateKey))
	log.Println("PUBLICKEY   :", status(config.PublicKey))
	log.Println("MONGOSTRING :", status(config.MongoString))

	srv := &http.Server{
		Addr:              config.Port,
		Handler:           http.HandlerFunc(routes.HandleRoutes),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	log.Println("Listening on", config.Port)
	log.Fatal(srv.ListenAndServe())
}

// status hanya menampilkan apakah env terisi, tanpa membocorkan nilainya ke log.
func status(v string) string {
	if v == "" {
		return "KOSONG"
	}
	return "terisi"
}
