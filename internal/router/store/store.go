package routerstore

import "github.com/go-chi/chi/v5"

var r *chi.Mux

func SetRouter(router *chi.Mux) {
    r = router
}

func GetRouter() *chi.Mux {
    return r
}
