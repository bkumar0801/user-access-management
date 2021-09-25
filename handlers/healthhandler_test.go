package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-access-management/mock"

	"github.com/gin-gonic/gin"
)

func TestHealthCheckHandler(t *testing.T) {
	server := gin.New()

	mockService := &mock.ServiceMock{}
	mockTokenManager := &mock.TokenManagerMock{}
	handler := NewUserAccessManagementHandler(mockService, mockTokenManager)

	server.Handle(http.MethodGet, "/health", handler.HealthHandler)
	httpServer := httptest.NewServer(server)

	cases := map[string]struct {
		dbErr mock.ErrMock
		code  int
		want  gin.H
	}{
		"health API responded with DB status": {
			dbErr: mock.OK,
			code:  http.StatusOK,
			want: gin.H{
				"status": "alive",
				"db":     "connected",
			},
		},
		"health API responded with DB error": {
			dbErr: mock.DBConnectionError,
			code:  http.StatusOK,
			want: gin.H{
				"status": "alive",
				"db":     "mock DB connection error",
			},
		},
		"health API responded with DB timeout": {
			dbErr: mock.DBConnectionTimeout,
			code:  http.StatusOK,
			want: gin.H{
				"status": "alive",
				"db":     "false",
			},
		},
	}

	gin.SetMode(gin.TestMode)
	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			mockService.Err = v.dbErr
			client := http.Client{}
			requestURL := httpServer.URL + "/health"
			req, err := http.NewRequest(http.MethodGet, requestURL, nil)
			if err != nil {
				t.Error("unexpected error: ", err)
			}
			req.Header.Set("authorization", "bearer token")

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

			if status := res.StatusCode; status != v.code {
				t.Errorf("handler returned wrong status code: \ngot %v\nwant %v\n", status, v.code)
			}

			if fmt.Sprint(v.want) != fmt.Sprint(got) {
				t.Errorf("handler returned unexpected body: \ngot %v\nwant %v\n", got, v.want)
			}
		})
	}

}
