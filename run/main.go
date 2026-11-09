package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/routes"
)

func main() {
	config.SetEnv()
	log.Println("Cek Environment Variabel")
	log.Println("LISTEN      :", config.Port)
	log.Println("PRIVATEKEY  :", status(config.PrivateKey))
	log.Println("PUBLICKEY   :", status(config.PublicKey))
	log.Println("MONGOSTRING :", status(config.MongoString))
	log.Println("DEMO LOGIN  :", demoStatus())

	srv := &http.Server{
		Addr:              config.Port,
		Handler:           http.HandlerFunc(routes.HandleRoutes),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Graceful shutdown: saat SIGINT/SIGTERM (Ctrl+C, systemctl stop/restart),
	// request yang sedang berjalan diberi waktu selesai sebelum server berhenti.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Listening on", config.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Shutdown error:", err)
	}
}

// status hanya menampilkan apakah env terisi, tanpa membocorkan nilainya ke log.
func status(v string) string {
	if v == "" {
		return "KOSONG"
	}
	return "terisi"
}

func demoStatus() string {
	if config.DemoUser == "" || config.DemoPassword == "" {
		return "mati (DEMO_USER/DEMO_PASSWORD kosong)"
	}
	return "aktif"
}
