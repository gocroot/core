package config

import (
	"net/http"
	"strings"
	"sync"
)

// defaultOrigins dipakai jika env ALLOWED_ORIGINS kosong.
var defaultOrigins = []string{
	"https://jscroot.if.co.id",
	"https://naskah.bukupedia.co.id",
	"https://bukupedia.co.id",
}

var (
	originsMu      sync.RWMutex
	allowedOrigins = toSet(defaultOrigins)
)

func toSet(origins []string) map[string]bool {
	set := make(map[string]bool, len(origins))
	for _, o := range origins {
		if o = strings.TrimSpace(o); o != "" {
			set[o] = true
		}
	}
	return set
}

// Fungsi untuk memeriksa apakah origin diizinkan (O(1) complexity)
func isAllowedOrigin(origin string) bool {
	originsMu.RLock()
	defer originsMu.RUnlock()
	return allowedOrigins[origin]
}

// Helper function untuk set common CORS headers
func setCORSHeaders(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", origin)
}

// SetAccessControlHeaders mengatur header CORS.
// Mengembalikan true jika request preflight (OPTIONS) sudah dijawab.
// Request tanpa header Origin (server-to-server, curl) tetap diproses tanpa header CORS.
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
	// Respons berbeda per Origin, jadi cache di tengah jalan harus membedakannya untuk semua respons.
	w.Header().Add("Vary", "Origin")
	origin := r.Header.Get("Origin")

	// Jika origin kosong atau tidak diizinkan, tidak ada header CORS
	if origin == "" || !isAllowedOrigin(origin) {
		return false
	}

	// Handle preflight request (OPTIONS)
	if r.Method == http.MethodOptions {
		setCORSHeaders(w, origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Token, Login, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return true
	}

	// Set CORS headers untuk request biasa
	setCORSHeaders(w, origin)
	return false
}

// SetAllowedOrigins mengganti seluruh daftar origin yang diizinkan.
func SetAllowedOrigins(origins []string) {
	set := toSet(origins)
	originsMu.Lock()
	allowedOrigins = set
	originsMu.Unlock()
}

// AddAllowedOrigin menambah origin baru (aman dipanggil saat server berjalan).
func AddAllowedOrigin(origin string) {
	originsMu.Lock()
	allowedOrigins[origin] = true
	originsMu.Unlock()
}

// RemoveAllowedOrigin menghapus origin (aman dipanggil saat server berjalan).
func RemoveAllowedOrigin(origin string) {
	originsMu.Lock()
	delete(allowedOrigins, origin)
	originsMu.Unlock()
}
