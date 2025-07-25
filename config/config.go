package config

import (
	"os"
	"sync"
)

var (
	once        sync.Once
	MongoString string
)

// SetEnv dengan protection untuk multiple calls
func SetEnv() {
	once.Do(func() {
		// Load environment variables
		// Set default configurations
		// Initialize global settings
		MongoString = os.Getenv("MONGOSTRING")
	})
}
