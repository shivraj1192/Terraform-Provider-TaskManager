package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockTaskManagerClient is a mock implementation of the TaskManagerClient
type MockTaskManagerClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
	Responses  map[string]MockResponse
}

// MockResponse represents a mock API response
type MockResponse struct {
	StatusCode int
	Body       interface{}
}

// NewMockClient creates a new MockTaskManagerClient
func NewMockClient() *MockTaskManagerClient {
	return &MockTaskManagerClient{
		BaseURL:    "http://mock-api/",
		Token:      "mock-token",
		HTTPClient: &http.Client{},
		Responses:  make(map[string]MockResponse),
	}
}

// AddMockResponse adds a mock response for a specific endpoint and method
func (c *MockTaskManagerClient) AddMockResponse(method, endpoint string, statusCode int, body interface{}) {
	key := method + ":" + endpoint
	c.Responses[key] = MockResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}

// Get is a mock implementation of the Get method
func (c *MockTaskManagerClient) Get(endpoint string, result interface{}) error {
	key := "GET:" + endpoint
	response, ok := c.Responses[key]
	if !ok {
		return fmt.Errorf("no mock response for GET %s", endpoint)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d)", response.StatusCode)
	}

	// Convert the mock response body to the result type
	bytes, err := json.Marshal(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, result)
}

// Post is a mock implementation of the Post method
func (c *MockTaskManagerClient) Post(endpoint string, body, result interface{}) error {
	key := "POST:" + endpoint
	response, ok := c.Responses[key]
	if !ok {
		return fmt.Errorf("no mock response for POST %s", endpoint)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d)", response.StatusCode)
	}

	// Convert the mock response body to the result type
	if result != nil {
		bytes, err := json.Marshal(response.Body)
		if err != nil {
			return err
		}

		return json.Unmarshal(bytes, result)
	}

	return nil
}

// Put is a mock implementation of the Put method
func (c *MockTaskManagerClient) Put(endpoint string, body, result interface{}) error {
	key := "PUT:" + endpoint
	response, ok := c.Responses[key]
	if !ok {
		return fmt.Errorf("no mock response for PUT %s", endpoint)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d)", response.StatusCode)
	}

	// Convert the mock response body to the result type
	if result != nil {
		bytes, err := json.Marshal(response.Body)
		if err != nil {
			return err
		}

		return json.Unmarshal(bytes, result)
	}

	return nil
}

// Delete is a mock implementation of the Delete method
func (c *MockTaskManagerClient) Delete(endpoint string) error {
	key := "DELETE:" + endpoint
	response, ok := c.Responses[key]
	if !ok {
		return fmt.Errorf("no mock response for DELETE %s", endpoint)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("API error (status %d)", response.StatusCode)
	}

	return nil
}

// CreateMockServer creates a test server that mocks the TaskManager API
func CreateMockServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the Authorization header is set
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized",
			})
			return
		}

		// Handle different endpoints
		switch r.URL.Path {
		case "/api/users":
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"id":    1,
						"uname": "testuser",
						"name":  "Test User",
						"email": "test@example.com",
					},
				})
			} else if r.Method == "POST" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":    1,
					"uname": "testuser",
					"name":  "Test User",
					"email": "test@example.com",
				})
			}
		case "/api/users/1":
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":    1,
					"uname": "testuser",
					"name":  "Test User",
					"email": "test@example.com",
				})
			} else if r.Method == "PUT" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":    1,
					"uname": "testuser",
					"name":  "Updated User",
					"email": "updated@example.com",
				})
			} else if r.Method == "DELETE" {
				w.WriteHeader(http.StatusNoContent)
			}
		case "/api/teams":
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"id":          1,
						"name":        "Test Team",
						"description": "A team for testing",
						"members":     []int{1},
					},
				})
			} else if r.Method == "POST" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":          1,
					"name":        "Test Team",
					"description": "A team for testing",
					"members":     []int{1},
				})
			}
		case "/api/tasks":
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{
						"id":          1,
						"title":       "Test Task",
						"description": "A task for testing",
						"status":      "To Do",
						"priority":    "Medium",
						"creator_id":  1,
						"team_id":     1,
						"assignees":   []int{1},
					},
				})
			} else if r.Method == "POST" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":          1,
					"title":       "Test Task",
					"description": "A task for testing",
					"status":      "To Do",
					"priority":    "Medium",
					"creator_id":  1,
					"team_id":     1,
					"assignees":   []int{1},
				})
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Not found",
			})
		}
	}))
}