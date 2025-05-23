package home

import (
	"GOrion/config/handlers"
	"GOrion/internal/router/store"
)

func HomeRoutes() {
	r := routerstore.GetRouter()

	// Your Routes here
	r.Get("/home", handlers.HomeHandler)
}