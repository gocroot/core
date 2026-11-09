package config

import (
	"net"
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

	// DemoUser dan DemoPassword mengaktifkan login demo. Jika salah satu kosong, login dimatikan.
	DemoUser     string
	DemoPassword string
)

// SetEnv dengan protection untuk multiple calls
func SetEnv() {
	once.Do(func() {
		Port = listenAddr(os.Getenv("IP"), os.Getenv("PORT"))
		MongoString = os.Getenv("MONGOSTRING")
		PrivateKey = os.Getenv("PRIVATEKEY")
		PublicKey = os.Getenv("PUBLICKEY")
		DemoUser = os.Getenv("DEMO_USER")
		DemoPassword = os.Getenv("DEMO_PASSWORD")
		if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
			SetAllowedOrigins(strings.Split(origins, ","))
		}
	})
}

// listenAddr menyusun alamat listen dari env IP dan PORT.
//
//	IP kosong, PORT kosong   -> ":3000" (semua interface)
//	IP kosong, PORT "8080"   -> ":8080"
//	IP "127.1.2.3"           -> "127.1.2.3:8080"
//	IP "::1" atau "[::1]"    -> "[::1]:8080"
//
// IP dibutuhkan di hosting yang memberi alamat khusus per akun, misalnya alwaysdata.
func listenAddr(ip, port string) string {
	ip = strings.Trim(strings.TrimSpace(ip), "[]")
	port = strings.TrimPrefix(strings.TrimSpace(port), ":")
	if port == "" {
		port = DefaultPort
	}
	return net.JoinHostPort(ip, port)
}
