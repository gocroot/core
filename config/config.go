package config

import (
	"os"
	"sync"
)

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
		// Load environment variables
		// Set default configurations
		// Initialize global settings
		Port = ":3000"
		MongoString = os.Getenv("MONGOSTRING")
		PrivateKey = os.Getenv("PRIVATEKEY")
		PublicKey = os.Getenv("PUBLICKEY")
	})
}
