package controller

import (
	"net/http"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/models"
	"github.com/whatsauth/watoken"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var loginauth models.LoginAuth
	if err := DecodeJSON(w, r, &loginauth); err != nil {
		WriteJSON(w, http.StatusBadRequest, models.Response{Status: "Illegal Request", Message: "Bad Request"})
		return
	}
	//dummy login, jika auth benar user admin password admin123
	if loginauth.UserName != "admin" || loginauth.Password != "admin123" {
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

func GetUserName(w http.ResponseWriter, r *http.Request) {
	useridstring, err := watoken.DecodeGetId(config.PublicKey, r.Header.Get("Token"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, models.Response{Status: "Illegal Request", Message: "Bad Request"})
		return
	}
	if useridstring == "" {
		WriteJSON(w, http.StatusNotFound, models.Response{Status: "Not Found", Message: "User Not Found"})
		return
	}
	WriteJSON(w, http.StatusOK, models.Response{Status: "Success", Message: "Token Valid", Username: useridstring})
}
