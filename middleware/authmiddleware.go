package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware ...
type AuthMiddleware struct {
	baseURL string
}

// NewAuthMiddleware is a constructor which create an object of AuthMiddleware
func NewAuthMiddleware(baseURL string) *AuthMiddleware {
	return &AuthMiddleware{
		baseURL: baseURL,
	}
}

// DoAuthenticate responds with unauthorized if header auth token is not valid
func (a *AuthMiddleware) DoAuthenticate(c *gin.Context) {
	userID, ok := c.Params.Get("userID")
	bearerToken := c.Request.Header.Get("authorization")
	if !ok || len(bearerToken) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "authentication token was not found in the request",
		})
		return
	}

	url := fmt.Sprintf("%s/users/%s/validatetoken", a.baseURL, userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", bearerToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "token verification failed",
		})
		return
	}
	c.Next()
}
