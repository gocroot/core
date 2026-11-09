package controller

import (
	"crypto/subtle"
	"net/http"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/models"
	"github.com/whatsauth/watoken"
)

// Login adalah contoh login demo. Kredensial diambil dari env DEMO_USER dan DEMO_PASSWORD.
// Jika salah satu kosong, login dimatikan. Ganti dengan pengecekan user di database untuk produksi.
func Login(w http.ResponseWriter, r *http.Request) {
	if config.DemoUser == "" || config.DemoPassword == "" {
		WriteJSON(w, http.StatusServiceUnavailable, models.Response{Status: "Error", Message: "Login belum dikonfigurasi"})
		return
	}
	var loginauth models.LoginAuth
	if err := DecodeJSON(w, r, &loginauth); err != nil {
		WriteJSON(w, http.StatusBadRequest, models.Response{Status: "Illegal Request", Message: "Bad Request"})
		return
	}
	// Bandingkan username dan password dengan waktu konstan, dan keduanya selalu dicek,
	// supaya lama respons tidak membocorkan bagian mana yang salah.
	userOK := subtle.ConstantTimeCompare([]byte(loginauth.UserName), []byte(config.DemoUser))
	passOK := subtle.ConstantTimeCompare([]byte(loginauth.Password), []byte(config.DemoPassword))
	if userOK&passOK != 1 {
		WriteJSON(w, http.StatusNotFound, models.Response{Status: "Not Found", Message: "User Not Found"})
		return
	}
	tokenstring, err := watoken.Encode(loginauth.UserName, config.PrivateKey)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, models.Response{Status: "Error", Message: "Gagal membuat token"})
		return
	}
	WriteJSON(w, http.StatusOK, models.Response{Status: "Success", Message: "Login successful", Token: tokenstring})
}

// GetUserName memvalidasi header Token. Token palsu, rusak, atau kedaluwarsa dibalas 401
// supaya frontend tahu user harus login ulang.
func GetUserName(w http.ResponseWriter, r *http.Request) {
	useridstring, err := watoken.DecodeGetId(config.PublicKey, r.Header.Get("Token"))
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, models.Response{Status: "Unauthorized", Message: "Token tidak valid atau kedaluwarsa"})
		return
	}
	if useridstring == "" {
		WriteJSON(w, http.StatusNotFound, models.Response{Status: "Not Found", Message: "User Not Found"})
		return
	}
	WriteJSON(w, http.StatusOK, models.Response{Status: "Success", Message: "Token Valid", Username: useridstring})
}
