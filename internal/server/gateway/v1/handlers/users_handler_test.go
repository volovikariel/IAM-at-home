package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/handlers"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/models"
)

func TestUserHandlerCreateUser(t *testing.T) {
	testCases := []struct {
		name           string
		userStore      *models.InMemoryUserStore
		request        mockRequest
		response       mockResponse
		expectedStatus int
	}{
		{
			name: "Create User Success",
			request: mockRequest{
				method:      "POST",
				path:        "/v1/users",
				contentType: "application/json",
				body:        `{"username":"newuser","password":"newpassword"}`, // Assuming JSON payload
			},
			response: mockResponse{
				contentType: "application/json",
				body:        `{"username":"newuser","rel":{"self":"/v1/users/newuser"}}`,
			},
			expectedStatus: http.StatusCreated, // Assuming 201 Created is the success status code
		},
		{
			name: "Create User Method Not Allowed",
			request: mockRequest{
				method: "GET", // Invalid method for creating a user
				path:   "/v1/users",
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "Create User Method With wrong Content-Type",
			request: mockRequest{
				method:      "GET", // Invalid method for creating a user
				path:        "/v1/users",
				contentType: "application/xml",
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "Create User Method With wrong fields in application/json body",
			request: mockRequest{
				method:      "POST",
				path:        "/v1/users",
				contentType: "application/json",
				body:        `{"username":"newuser"}`, // Missing password field
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Create User Method fails because of duplicate user",
			request: mockRequest{
				method:      "POST",
				path:        "/v1/users",
				contentType: "application/json",
				body:        `{"username":"newuser", "password":"newpassword"}`, // Assuming JSON payload
			},
			userStore: models.NewInMemoryUserStore([]models.User{
				{
					Name:     "newuser",
					Password: "newpassword",
				},
			}),
			expectedStatus: http.StatusConflict,
		},
		{
			name: "Create User Method fails when body of request is not correct for the specified content type",
			request: mockRequest{
				method:      "POST",
				path:        "/v1/users",
				contentType: "application/json",
				body:        "username:newuser", // Not valid JSON
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		tc := tc // Prevent race condition given parallel test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup userStore and sessionStore as needed...
			userStore := models.InMemoryUserStore{}
			if tc.userStore != nil {
				userStore = *tc.userStore
			}
			sessionStore := models.InMemorySessionStore{}
			handler := handlers.NewV1Handler(&userStore, &sessionStore)

			// Convert the string body to a reader, which could be nil if no body is required
			bodyReader := strings.NewReader(tc.request.body)
			req, err := http.NewRequest(tc.request.method, tc.request.path, bodyReader)
			req.Header.Set("Content-Type", tc.request.contentType)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Check for the expected status code
			if status := rr.Result().StatusCode; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
			// if we don't expect a body, we don't need to check it (as bodies sometimes returned are just the error message, and we already check that with the status code)
			if tc.response.body != "" && rr.Body.String() != tc.response.body {
				t.Errorf("handler returned wrong body: got %v want %v", rr.Body.String(), tc.response.body)
			}
		})
	}
}

// TODO: Handle 200
// TODO: Handle 404
func TestUserHandlerGetUser(t *testing.T) {
	testCases := []struct {
		name           string
		userStore      *models.InMemoryUserStore
		request        mockRequest
		response       mockResponse
		expectedStatus int
	}{
		{
			name: "Get User Success",
			request: mockRequest{
				method: "GET",
				path:   "/v1/users/newuser",
			},
			response: mockResponse{
				contentType: "application/json",
				body:        `{"username":"newuser", "password":"newpassword", "rel":{"self":"/v1/users/newuser"}}`,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Get User Method Not Allowed",
			request: mockRequest{
				method: "POST",
				path:   "/v1/users/newuser",
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tc := range testCases {
		tc := tc // Prevent race condition given parallel test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup userStore and sessionStore as needed...
			userStore := models.InMemoryUserStore{}
			if tc.userStore != nil {
				userStore = *tc.userStore
			}
			sessionStore := models.InMemorySessionStore{}
			handler := handlers.NewV1Handler(&userStore, &sessionStore)

			// Convert the string body to a reader, which could be nil if no body is required
			bodyReader := strings.NewReader(tc.request.body)
			req, err := http.NewRequest(tc.request.method, tc.request.path, bodyReader)
			req.Header.Set("Content-Type", tc.request.contentType)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Check for the expected status code
			if status := rr.Result().StatusCode; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
			// if we don't expect a body, we don't need to check it (as bodies sometimes returned are just the error message, and we already check that with the status code)
			if tc.response.body != "" && rr.Body.String() != tc.response.body {
				t.Errorf("handler returned wrong body: got %v want %v", rr.Body.String(), tc.response.body)
			}
		})
	}
}

// TODO: Handle 200 (+body)
// TODO: Handle 400
// TODO: Handle 401
// TODO: Handle 403
// TODO: Handle 404
// TODO: Handle 415
// TODO: Handle 422
func TestUserHandlerUpdateUser(t *testing.T) {
	testCases := []struct {
		name           string
		userStore      *models.InMemoryUserStore
		request        mockRequest
		response       mockResponse
		expectedStatus int
	}{
		{
			name: "Update User Success",
			request: mockRequest{
				method: "PATCH",
				path:   "/v1/users/newuser",
			},
			response: mockResponse{
				contentType: "application/json",
				body:        `{"username":"newuser", "password":"newpassword", "rel":{"self":"/v1/users/newuser"}}`,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Update User Method Not Allowed",
			request: mockRequest{
				method: "POST",
				path:   "/v1/users/newuser",
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tc := range testCases {
		tc := tc // Prevent race condition given parallel test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup userStore and sessionStore as needed...
			userStore := models.InMemoryUserStore{}
			if tc.userStore != nil {
				userStore = *tc.userStore
			}
			sessionStore := models.InMemorySessionStore{}
			handler := handlers.NewV1Handler(&userStore, &sessionStore)

			// Convert the string body to a reader, which could be nil if no body is required
			bodyReader := strings.NewReader(tc.request.body)
			req, err := http.NewRequest(tc.request.method, tc.request.path, bodyReader)
			req.Header.Set("Content-Type", tc.request.contentType)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Check for the expected status code
			if status := rr.Result().StatusCode; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
			// if we don't expect a body, we don't need to check it (as bodies sometimes returned are just the error message, and we already check that with the status code)
			if tc.response.body != "" && rr.Body.String() != tc.response.body {
				t.Errorf("handler returned wrong body: got %v want %v", rr.Body.String(), tc.response.body)
			}
		})
	}
}

// TODO: Handle 204
// TODO: Handle 400
// TODO: Handle 401 (+body)
// TODO: Handle 403
// TODO: Handle 404
// TODO: Handle 422
func TestUserHandlerDeleteUser(t *testing.T) {
	testCases := []struct {
		name           string
		userStore      *models.InMemoryUserStore
		request        mockRequest
		response       mockResponse
		expectedStatus int
	}{
		{
			name: "Delete User Success",
			request: mockRequest{
				method:      "DELETE",
				path:        "/v1/users/newuser",
				contentType: "application/json",
			},
			response: mockResponse{
				contentType: "application/json",
				body:        `{"username":"newuser", "rel":{"self":"/v1/users/newuser"}}`,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Delete User Method Not Allowed",
			request: mockRequest{
				method: "POST",
				path:   "/v1/users/newuser",
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "Delete User Not Found",
			request: mockRequest{
				method:      "DELETE",
				path:        "/v1/users/newuser",
				contentType: "application/json",
			},
			expectedStatus: http.StatusNotFound,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Setup userStore and sessionStore as needed...
			userStore := models.InMemoryUserStore{}
			if tc.userStore != nil {
				userStore = *tc.userStore
			}
			sessionStore := models.InMemorySessionStore{}
			handler := handlers.NewV1Handler(&userStore, &sessionStore)

			// Convert the string body to a reader, which could be nil if no body is required
			bodyReader := strings.NewReader(tc.request.body)
			req, err := http.NewRequest(tc.request.method, tc.request.path, bodyReader)
			req.Header.Set("Content-Type", tc.request.contentType)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Check for the expected status code
			if status := rr.Result().StatusCode; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
			// if we don't expect a body, we don't need to check it (as bodies sometimes returned are just the error message, and we already check that with the status code)
			if tc.response.body != "" && rr.Body.String() != tc.response.body {
				t.Errorf("handler returned wrong body: got %v want %v", rr.Body.String(), tc.response.body)
			}
		})
	}
}
