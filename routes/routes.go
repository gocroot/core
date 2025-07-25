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

	var method, path string = r.Method, r.URL.Path
	switch {
	case method == "GET" && path == "/":
		controller.GetHome(w, r)

	default:
		controller.NotFound(w, r)
	}
}
