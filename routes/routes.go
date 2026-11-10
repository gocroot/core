package routes

import (
	"net/http"

	"github.com/gocroot/core/config"
	"github.com/gocroot/core/controller"
)

func HandleRoutes(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}

	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/":
			controller.GetHome(w, r)
		case "/login":
			controller.GetUserName(w, r)
		case "/version":
			controller.GetVersion(w, r)
		default:
			// Route dengan path parameter ditaruh di default, setelah semua path statis.
			// Contoh /data/user/{id} (ambil nilainya di controller dengan PathParam yang sama):
			//
			//	if _, ok := controller.PathParam(r.URL.Path, "/data/user/"); ok {
			//		controller.GetUserByID(w, r)
			//		return
			//	}
			controller.NotFound(w, r)
		}

	case http.MethodPost:
		switch r.URL.Path {
		case "/login":
			controller.Login(w, r)
		default:
			controller.NotFound(w, r)
		}

	default:
		controller.MethodNotAllowed(w, r)
	}
}
