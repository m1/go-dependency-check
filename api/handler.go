package api

import (
	"github.com/go-chi/chi"
)

// Handler is the interface used by all the handlers to
// pass back all of their routes
type Handler interface {
	GetRoutes() *chi.Mux
}
