package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/routes"
)

func main() {
	config.SetEnv()
	log.Println("Cek Environment Variabel")
	log.Println("VERSION     :", config.Version)
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

	/*
		Berhenti dengan rapi saat supervisor hosting mengirim sinyal.

		Tanpa penanganan sinyal, proses berhenti seketika di tengah request
		yang sedang berjalan, sehingga data bisa terputus separuh tertulis.
		Sebaliknya, proses yang tidak pernah benar-benar keluar akan terus
		melayani dengan kode lama meski binary baru sudah terunggah:
		"restart" dilaporkan berhasil, tetapi tidak ada yang berubah.

		Penutupan berbatas waktu menutup kedua celah itu: request yang sedang
		jalan diberi kesempatan selesai, dan prosesnya dipastikan keluar.

		Daftar sinyal ada di sinyal_unix.go: SIGINT, SIGTERM, dan SIGHUP,
		SIGUSR1, SIGUSR2 (pilihan sinyal "Hot restart" di panel hosting
		seperti alwaysdata).
	*/
	berhenti := make(chan os.Signal, 1)
	signal.Notify(berhenti, sinyalBerhenti...)

	gagal := make(chan error, 1)
	go func() {
		log.Println("Listening on", config.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			gagal <- err
		}
	}()

	select {
	case err := <-gagal:
		// Misalnya port sudah dipakai: keluar dengan kode error supaya supervisor tahu.
		log.Fatal(err)
	case s := <-berhenti:
		log.Printf("sinyal %v diterima, menutup dengan rapi (maks %s)", s, masaTutup)
	}
	signal.Stop(berhenti)

	ctx, cancel := context.WithTimeout(context.Background(), masaTutup)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Batas waktu habis: koneksi yang tersisa ditutup paksa, proses tetap keluar.
		log.Println("penutupan tidak tuntas:", err)
		srv.Close()
	}
}

// masaTutup adalah batas waktu request yang sedang berjalan untuk selesai saat
// server dihentikan. Harus lebih pendek dari jeda sebelum supervisor mengirim
// SIGKILL (systemd 90 detik, Docker dan Cloud Run 10 detik secara default).
const masaTutup = 8 * time.Second

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
