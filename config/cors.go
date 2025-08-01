package config

import (
	"net/http"
)

// Menggunakan map untuk lookup O(1) - lebih efisien
var allowedOrigins = map[string]bool{
	"https://jscroot.if.co.id":       true,
	"https://naskah.bukupedia.co.id": true,
	"https://bukupedia.co.id":        true,
}

// Fungsi untuk memeriksa apakah origin diizinkan (O(1) complexity)
func isAllowedOrigin(origin string) bool {
	return allowedOrigins[origin]
}

// Helper function untuk set common CORS headers
func setCORSHeaders(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", origin)
}

// Fungsi untuk mengatur header CORS
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// Jika origin kosong atau tidak diizinkan, reject
	if origin == "" || !isAllowedOrigin(origin) {
		return false
	}

	// Handle preflight request (OPTIONS)
	if r.Method == http.MethodOptions {
		setCORSHeaders(w, origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Login, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return true
	}

	// Set CORS headers untuk request biasa
	setCORSHeaders(w, origin)
	return false
}

// Optional: Fungsi untuk menambah origin baru secara dinamis
func AddAllowedOrigin(origin string) {
	allowedOrigins[origin] = true
}

// Optional: Fungsi untuk remove origin
func RemoveAllowedOrigin(origin string) {
	delete(allowedOrigins, origin)
}
