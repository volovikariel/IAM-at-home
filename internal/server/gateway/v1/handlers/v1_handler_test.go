package handlers_test

// TODO: See backlog

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/handlers"
// 	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/models"
// )

// // TODO: Make these explicitely shared (internal/testing?)
// // or perhaps, have a request struct actually used in the code that we can use here
// type mockRequest struct {
// 	method      string
// 	path        string
// 	contentType string
// 	body        string
// }

// type mockResponse struct {
// 	status      int
// 	contentType string
// 	body        string
// }

// func TestV1HandlerServeHTTP(t *testing.T) {
// 	testCases := []struct {
// 		name     string
// 		request  mockRequest
// 		response mockResponse
// 	}{
// 		{
// 			name: "Path not starting with forward slash should return not found",
// 			request: mockRequest{
// 				method: "POST",
// 				path:   "v1/users",
// 			},
// 			response: mockResponse{
// 				status: http.StatusNotFound,
// 			},
// 		},
// 		{
// 			name: "Trailing slash /v1/users/ should return not found",
// 			request: mockRequest{
// 				method: "POST",
// 				path:   "/v1/users/",
// 			},
// 			response: mockResponse{
// 				status: http.StatusNotFound,
// 			},
// 		},
// 		{
// 			name: "Invalid path /v1/userssss should return not found",
// 			request: mockRequest{
// 				method: "POST",
// 				path:   "/v1/userssss",
// 			},
// 			response: mockResponse{
// 				status: http.StatusNotFound,
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc // Prevent race condition given parallel test
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			userStore := models.InMemoryUserStore{}
// 			sessionStore := models.InMemorySessionStore{}
// 			handler := handlers.NewV1Handler(&userStore, &sessionStore)
// 			req, err := http.NewRequest(tc.request.method, tc.request.path, nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			rr := httptest.NewRecorder()
// 			handler.ServeHTTP(rr, req)
// 			if status := rr.Code; status != tc.response.status {
// 				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.response.status)
// 			}
// 		})
// 	}
// }

// func TestV1HandlerServeHTTPDoesntReturnNotFound(t *testing.T) {
// 	testCases := []struct {
// 		name    string
// 		request mockRequest
// 	}{
// 		{
// 			name: "A correct path should not return not found",
// 			request: mockRequest{
// 				method: "POST",
// 				path:   "/v1/users",
// 			},
// 		},
// 	}
// 	for _, tc := range testCases {
// 		tc := tc // Prevent race condition given parallel test
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			userStore := models.InMemoryUserStore{}
// 			sessionStore := models.InMemorySessionStore{}
// 			handler := handlers.NewV1Handler(&userStore, &sessionStore)
// 			req, err := http.NewRequest(tc.request.method, tc.request.path, nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			rr := httptest.NewRecorder()
// 			handler.ServeHTTP(rr, req)
// 			if rr.Result().StatusCode == http.StatusNotFound {
// 				t.Errorf("handler returned undesireable status code: %v", http.StatusNotFound)
// 			}
// 		})
// 	}
// }
