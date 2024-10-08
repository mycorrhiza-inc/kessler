package crud

import (
	"github.com/gorilla/mux"
)

func defineCrudRoutes(router *mux.Router) {
	s := router.PathPrefix("/crud").Subrouter()
}
