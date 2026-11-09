package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/core/models"
)

// maxBodyBytes membatasi ukuran body request JSON (1 MiB).
const maxBodyBytes = 1 << 20

// WriteJSON menulis header Content-Type, status code, lalu body JSON.
// Respons API tidak boleh disimpan cache (bisa berisi token) dan tidak boleh di-sniff sebagai tipe lain.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

// DecodeJSON membaca body JSON ke dst dengan batas ukuran maxBodyBytes.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	return json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(dst)
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, models.Response{Status: "success", Message: "Welcome"})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusNotFound, models.Response{Status: "Not Found", Message: "Page Not Found"})
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusMethodNotAllowed, models.Response{Status: "Method Not Allowed", Message: "Method " + r.Method + " Not Allowed"})
}
