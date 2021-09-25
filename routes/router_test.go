package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"user-access-management/handlers"
	"user-access-management/middleware"
	"user-access-management/mock"

	"github.com/gin-gonic/gin"
)

func TestNewRoutes(t *testing.T) {
	mockService := &mock.ServiceMock{
		Err: mock.OK,
	}
	mockTokenManager := &mock.TokenManagerMock{}
	userAccessManagementHandler := handlers.NewUserAccessManagementHandler(mockService, mockTokenManager)
	got := NewRoutes(userAccessManagementHandler)
	expectedRoutes := []string{
		"/health",
	}
	t.Run("all routes are present", func(t *testing.T) {
		for i, v := range got {
			if v.Pattern != expectedRoutes[i] {
				t.Errorf("pattern expectations mismatched: \n want: %v \n got: %v", v.Pattern, v)
			}
		}
	})
}

func TestAttachRoutes(t *testing.T) {
	mockService := &mock.ServiceMock{
		Err: mock.OK,
	}
	mockTokenManager := &mock.TokenManagerMock{}
	userAccessManagementHandler := handlers.NewUserAccessManagementHandler(mockService, mockTokenManager)
	routes := NewRoutes(userAccessManagementHandler)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	}))
	defer ts.Close()

	gin.SetMode(gin.TestMode)

	server := gin.New()
	authMiddelware := middleware.NewAuthMiddleware(ts.URL)
	httpServer := httptest.NewServer(server)

	AttachRoutes(server, routes, authMiddelware)

	cases := map[string]struct {
		pattern string
		method  string
		token   string
		body    *bytes.Buffer
		want    gin.H
		code    int
	}{
		"health endpoint is called": {
			pattern: "/health",
			method:  http.MethodGet,
			token:   "user t@t.com seller",
			body:    nil,
			want: gin.H{
				"status": "alive",
				"db":     "connected",
			},
			code: http.StatusOK,
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			requestURL := httpServer.URL + v.pattern
			var req *http.Request
			var err error
			if v.body != nil {
				req, err = http.NewRequest(v.method, requestURL, v.body)
				if err != nil {
					t.Error("unexpected error: ", err)
				}
			} else {
				req, err = http.NewRequest(v.method, requestURL, nil)
				if err != nil {
					t.Error("unexpected error: ", err)
				}
			}

			req.Header.Set("authorization", v.token)

			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Error("unexpected error: ", err)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error("unexpected error: ", err)
			}

			if !reflect.DeepEqual(v.want, gin.H{}) {
				var got gin.H
				err = json.Unmarshal(body, &got)
				if err != nil {
					t.Fatal(err)
				}

				if fmt.Sprint(v.want) != fmt.Sprint(got) {
					t.Errorf("attached router returned unexpected body: \ngot %v\nwant %v\n", got, v.want)
				}
			}

			if res.StatusCode != v.code {
				t.Errorf("status code mismatched: \n want: %v \ngot: %v", v.code, res.StatusCode)
			}

		})
	}

}
