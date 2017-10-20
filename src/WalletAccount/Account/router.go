package Account

import (
	"net/http"
	"gopkg.in/mux"

	//"WalletAccount/logger"
)

var controller = &Controller{Repository: Repository{}}
// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
//Routes defines the list of routes of our API
type Routes []Route
var routes = Routes{
	Route{
		"WalletAccount",
		"GET",
		"/",
		controller.Index,
	},
	Route{
		"AddAccount",
		"POST",
		"/",
		controller.AddAccount,
	},
	Route{
		"UpdateAccount",
		"PUT",
		"/",
		controller.UpdateAccount,
	},
	Route{
		"DeleteAccount",
		"DELETE",
		"/",
		controller.DeleteAccount,
	},
}
//NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)
		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
