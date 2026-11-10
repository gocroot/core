//go:build !unix

package main

import (
	"os"
	"syscall"
)

// sinyalBerhenti untuk OS non-Unix (Windows): SIGUSR1/SIGUSR2 tidak ada di sana.
// Daftar lengkap dan penjelasannya ada di sinyal_unix.go.
var sinyalBerhenti = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGHUP,
}
