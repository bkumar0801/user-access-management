package routes

import (
	"user-access-management/handlers"
	"user-access-management/middleware"

	"github.com/gin-gonic/gin"
)

/*
Route Structure of new routes
*/
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Protected   bool
	HandlerFunc gin.HandlerFunc
}

/*
Routes Array of all available routes
*/
type Routes []Route

/*
NewRoutes returns the list of available routes
*/
func NewRoutes(userAccessManagementHandler *handlers.UserAccessManagementHandler) Routes {
	return Routes{
		Route{
			"Health",
			"GET",
			"/health",
			false,
			userAccessManagementHandler.HealthHandler,
		},
	}
}

/*
AttachRoutes Attaches routes to the provided server
*/
func AttachRoutes(server *gin.Engine, routes Routes, authMiddleware *middleware.AuthMiddleware) {
	for _, route := range routes {
		if route.Protected {
			server.Handle(route.Method, route.Pattern, authMiddleware.DoAuthenticate, route.HandlerFunc)
		} else {
			server.Handle(route.Method, route.Pattern, route.HandlerFunc)
		}
	}
}
