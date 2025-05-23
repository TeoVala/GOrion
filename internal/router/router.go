package router

import (
	"GOrion/config/middleware"
	"GOrion/config/routes"
	"GOrion/internal/helpers/terminal"
	"GOrion/internal/router/store"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var r *chi.Mux

func SetupRouter() *chi.Mux {
	r = chi.NewRouter()
	
	routerstore.SetRouter(r)

	// Register Everything :)
	middleware.RegisterMiddlewares()
	routes.RegisterRoutes()
	
	return r
}

func GetAllRoutes() {

	var logMsg string
	logMsg += "\n--------------- Route list ---------------\n"
	// Print Route list
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logMsg += fmt.Sprintf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
	logMsg += "------------------------------------------\n"
	
	log.Printf(logMsg)

	terminal.CW(true, terminal.NWhite, "\n--------------- Route list ---------------\n")
	// Print Route list
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		terminal.CW(true, terminal.NBlue, "[%s]:", method,)
		terminal.CW(true, terminal.NGreen, " '%s' has ", route )
		terminal.CW(true, terminal.BRed, "%d middlewares\n", len(middlewares))
		return nil
	})
	terminal.CW(true, terminal.NWhite, "------------------------------------------\n")
}
