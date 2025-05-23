package middleware

import (
	"GOrion/config/middleware/requestlogger"
	routerstore "GOrion/internal/router/store"

	"github.com/go-chi/chi/v5/middleware"
)

func RegisterMiddlewares() {

	r := routerstore.GetRouter()

	// Add you middleware here
	r.Use(middleware.Recoverer)
	r.Use(requestlogger.RequestLogger)

}
