//go:build unix

package main

import (
	"os"
	"syscall"
)

// sinyalBerhenti adalah sinyal yang memicu penutupan dengan rapi.
//
//	SIGINT  : Ctrl+C saat go run ./run
//	SIGTERM : systemctl stop/restart, docker stop, Cloud Run, Render, Fly.io
//	SIGHUP, SIGUSR1, SIGUSR2 : pilihan sinyal "Hot restart" di panel hosting
//	          (misalnya alwaysdata). Tanpa penanganan, ketiganya membunuh
//	          proses Go seketika di tengah request yang sedang berjalan.
var sinyalBerhenti = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGHUP,
	syscall.SIGUSR1,
	syscall.SIGUSR2,
}
