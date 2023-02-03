package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
)

// REST query and parameter values
const (
	MethodGet = "GET"
)

// RegisterRoutes registers the cellarfees module REST routes.
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router, storeName string) {
	r := rest.WithHTTPDeprecationHeaders(rtr)
	r.HandleFunc(
		"/cellarfees/apy", QueryAPYRequestHandlerFn(storeName, clientCtx),
	).Methods(MethodGet)
}
