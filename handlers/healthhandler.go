package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
HealthHandler returns alive status
*/
func (u *UserAccessManagementHandler) HealthHandler(c *gin.Context) {
	ok, err := u.service.DBStatus()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
			"db":     err.Error(),
		})
		return
	}
	if ok {
		c.JSON(200, gin.H{
			"status": "alive",
			"db":     "connected",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "alive",
		"db":     "false",
	})
	return
}
