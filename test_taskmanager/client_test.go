package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"terraform-provider-taskmanager/taskmanager"
)

// TestNewClient tests the creation of a new TaskManagerClient
func TestNewClient(t *testing.T) {
	baseURL := "http://localhost:8080/"
	token := "test-token"

	client := taskmanager.NewClient(baseURL, token)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}
}

// setupMockServer creates a test server that mocks the TaskManager API
func setupMockServer(t *testing.T, path string, statusCode int, response interface{}) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request path matches the expected path
		if r.URL.Path != path {
			t.Errorf("Expected request to '%s', got '%s'", path, r.URL.Path)
		}

		// Check if the Authorization header is set
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", auth)
		}

		// Set response status code
		w.WriteHeader(statusCode)

		// Write response body if provided
		if response != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))

	return server
}

// TestClientGet tests the Get method of TaskManagerClient
func TestClientGet(t *testing.T) {
	// Setup mock response
	expectedResponse := map[string]interface{}{
		"id":    1,
		"name":  "Test User",
		"email": "test@example.com",
	}

	// Setup mock server
	server := setupMockServer(t, "/api/users/1", http.StatusOK, expectedResponse)
	defer server.Close()

	// Create client pointing to mock server
	client := taskmanager.NewClient(server.URL+"/", "test-token")

	// Test Get method
	var result map[string]interface{}
	err := client.Get("api/users/1", &result)

	// Check for errors
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check response data
	if result["id"].(float64) != 1 {
		t.Errorf("Expected id 1, got %v", result["id"])
	}

	if result["name"].(string) != "Test User" {
		t.Errorf("Expected name 'Test User', got %v", result["name"])
	}

	if result["email"].(string) != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %v", result["email"])
	}
}

// TestClientGetError tests error handling in the Get method
func TestClientGetError(t *testing.T) {
	// Setup mock server that returns an error
	server := setupMockServer(t, "/api/users/999", http.StatusNotFound, map[string]string{
		"error": "User not found",
	})
	defer server.Close()

	// Create client pointing to mock server
	client := taskmanager.NewClient(server.URL+"/", "test-token")

	// Test Get method with error response
	var result map[string]interface{}
	err := client.Get("api/users/999", &result)

	// Check for expected error
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}