package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/models"
	"github.com/whatsauth/watoken"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var loginauth models.LoginAuth
	err := json.NewDecoder(r.Body).Decode(&loginauth)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"status":"Illegal Request","message":"Bad Request"}`))
		return
	}
	//dummy login, jika auth benar user admin password admin123
	if loginauth.UserName == "admin" && loginauth.Password == "admin123" {
		w.WriteHeader(200)
		tokenstring, _ := watoken.Encode(loginauth.UserName, config.PrivateKey)
		response := fmt.Sprintf(`{"status":"Success","message":"Login successful","token":"%s"}`, tokenstring)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(404)
	w.Write([]byte(`{"status":"Not Found","message":"User Not FOund"}`))
}

func GetUserName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Token")
	useridstring, err := watoken.DecodeGetId(config.PublicKey, token)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(`{"status":"Illegal Request","message":"Bad Request"}`))
		return
	}
	//dummy login, jika auth benar user admin password admin123
	if useridstring != "" {
		w.WriteHeader(200)
		response := fmt.Sprintf(`{"status":"Success","message":"Token Valid","username":"%s"}`, useridstring)
		w.Write([]byte(response))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(`{"status":"Not Found","message":"User Not FOund"}`))
}
