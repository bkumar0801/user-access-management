package handlers

import (
	"reflect"
	"strings"
	"user-access-management/jwtparser"
	"user-access-management/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

//Handler is an interface for review CRUD operation
type Handler interface {
	HealthHandler(c *gin.Context)
}

/*
UserAccessManagementHandler handles the different handler types
*/
type UserAccessManagementHandler struct {
	service      service.Service
	tokenManager jwtparser.TokenManager
}

/*
NewUserAccessManagementHandler ...
*/
func NewUserAccessManagementHandler(service service.Service, tokenManager jwtparser.TokenManager) *UserAccessManagementHandler {
	return &UserAccessManagementHandler{
		service:      service,
		tokenManager: tokenManager,
	}
}
