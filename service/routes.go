package service

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"healthcheck",
		"GET",
		"/api/v1/rest/healthcheck",
		healthcheck,
	},
	Route{
		"doDelivery",
		"POST",
		"/api/rest/v1/deliver",
		deliveryHandler,
	},
}
