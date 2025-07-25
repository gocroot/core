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
	case "GET":
		switch r.URL.Path {
		case "/":
			controller.GetHome(w, r)
		case "/login":
			controller.GetProfile(w, r)
		}

	case "POST":
		switch r.URL.Path {
		case "/login":
			controller.Login(w, r)
		}
	default:
		controller.NotFound(w, r)
	}

}
