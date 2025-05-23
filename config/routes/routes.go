package routes

import (
	"GOrion/config/routes/home"
	routerstore "GOrion/internal/router/store"
	"net/http"
)

// RegisterRoutes defines all the core application routes
func RegisterRoutes() {
	r := routerstore.GetRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	home.HomeRoutes()
}