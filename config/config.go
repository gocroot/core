package config

import (
	"os"
	"strings"
	"sync"
)

// DefaultPort dipakai jika env PORT kosong.
const DefaultPort = "3000"

var (
	once        sync.Once
	Port        string
	MongoString string
	PrivateKey  string
	PublicKey   string
)

// SetEnv dengan protection untuk multiple calls
func SetEnv() {
	once.Do(func() {
		Port = portFromEnv(os.Getenv("PORT"))
		MongoString = os.Getenv("MONGOSTRING")
		PrivateKey = os.Getenv("PRIVATEKEY")
		PublicKey = os.Getenv("PUBLICKEY")
		if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
			SetAllowedOrigins(strings.Split(origins, ","))
		}
	})
}

// portFromEnv mengubah nilai env PORT ("8080" atau ":8080") menjadi alamat listen.
func portFromEnv(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		p = DefaultPort
	}
	if !strings.HasPrefix(p, ":") {
		p = ":" + p
	}
	return p
}
