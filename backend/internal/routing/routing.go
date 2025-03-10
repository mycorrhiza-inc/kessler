package routing

import "github.com/gorilla/mux"

type RouteDefinition struct {
	Router *mux.Router
	Prefix string
}
type Endpoint interface {
	DefineRoutes(RouteDefinition)
}
