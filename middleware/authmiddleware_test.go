package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDoAuthenticate(t *testing.T) {
	cases := map[string]struct {
		userID string
		token  string
		path   string
		status int
		resp   gin.H
	}{
		"user is authenticated": {
			userID: "test@test.com",
			token:  "bearer fakevalidtoken",
			path:   "",
			status: http.StatusOK,
			resp: gin.H{
				"message": "authentication successful",
			},
		},
		"received empty authentication token": {
			userID: "test@test.com",
			token:  "",
			path:   "",
			status: http.StatusUnauthorized,
			resp: gin.H{
				"message": "authentication token was not found in the request",
			},
		},
		"token service returned bad error code": {
			userID: "test@test.com",
			token:  "bearer fakevalidtoken",
			path:   "/error",
			status: http.StatusUnauthorized,
			resp: gin.H{
				"message": "token verification failed",
			},
		},
	}

	// Last fallbackHandler which will be executed after middleware
	fallbackHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "authentication successful",
		})
	}

	//setup mock token server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "error") {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}))
	defer ts.Close()

	gin.SetMode(gin.TestMode)

	server := gin.New()
	authMiddelware := NewAuthMiddleware(ts.URL)

	server.Handle(http.MethodGet, "/users/:userID", authMiddelware.DoAuthenticate, fallbackHandler)
	httpServer := httptest.NewServer(server)

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			authMiddelware.baseURL = authMiddelware.baseURL + v.path
			client := http.Client{}
			requestURL := httpServer.URL + "/users/" + v.userID
			req, err := http.NewRequest(http.MethodGet, requestURL, nil)
			if err != nil {
				t.Error("unexpected error: ", err)
			}
			req.Header.Set("authorization", v.token)

			res, err := client.Do(req)
			if err != nil {
				t.Error("unexpected error: ", err)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error("unexpected error: ", err)
			}

			var got gin.H
			err = json.Unmarshal(body, &got)
			if err != nil {
				t.Fatal(err)
			}

			if status := res.StatusCode; status != v.status {
				t.Errorf("middleware returned wrong status code: \ngot %v\nwant %v\n", status, v.status)
			}

			if fmt.Sprint(v.resp) != fmt.Sprint(got) {
				t.Errorf("middleware returned unexpected body: \ngot %v\nwant %v\n", got, v.resp)
			}
			authMiddelware.baseURL = ts.URL
		})
	}
}
