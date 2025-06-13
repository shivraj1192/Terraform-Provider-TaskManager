# TaskManager Provider Tests

This directory contains tests for the TaskManager Terraform provider. The tests are organized into several categories and are designed to ensure the provider works correctly with the TaskManager-Go API.

## Test Categories

- **Unit Tests**: Tests for individual functions and methods without requiring a running API
- **Acceptance Tests**: Tests that verify the provider works correctly with actual infrastructure
- **Provider Tests**: Tests that verify the provider configuration and initialization

## Prerequisites

Before running the tests, ensure you have:

1. Go 1.18 or later installed
2. A running instance of the TaskManager-Go API
3. Valid API credentials

## Environment Variables

The following environment variables are used by the tests:

- `TF_ACC`: Set to `1` to enable acceptance tests
- `TASKMANAGER_TOKEN`: Your API token for authentication
- `TASKMANAGER_BASE_URL`: The base URL of your TaskManager-Go API (default: `http://localhost:8080/`)

## Running the Tests

### Unit Tests

To run only the unit tests (which don't require a running API):

```sh
go test -v -short ./...
```

### Acceptance Tests

To run the acceptance tests (which require a running API):

```sh
# Set environment variables for the tests
export TF_ACC=1
export TASKMANAGER_TOKEN="your-api-token"
export TASKMANAGER_BASE_URL="http://localhost:8080/"

# Run the tests
go test -v ./...
```

On Windows (PowerShell):

```powershell
$env:TF_ACC=1
$env:TASKMANAGER_TOKEN="your-api-token"
$env:TASKMANAGER_BASE_URL="http://localhost:8080/"

go test -v ./...
```

### Testing Specific Files

To run tests in a specific file:

```sh
go test -v ./test_taskmanager -run TestProvider
```

This will run only the tests in the `provider_test.go` file.

## Test Structure

- `provider_test.go`: Tests for provider configuration and initialization
- `client_test.go`: Tests for the TaskManagerClient
- `resource_user_test.go`: Tests for the user resource
- `resource_team_test.go`: Tests for the team resource
- `resource_task_test.go`: Tests for the task resource
- `resource_comment_test.go`: Tests for the comment resource
- `resource_attachment_test.go`: Tests for the attachment resource
- `data_source_test.go`: Tests for all data sources
- `test_helper.go`: Helper functions for testing

## Important Notes on Test Implementation

### Provider Factories

In Terraform SDK v2, the `resource.TestCase` struct uses `ProviderFactories` instead of `Providers`. This is a key difference from SDK v1. The `ProviderFactories` field expects a map of functions that return a `*schema.Provider` and an error.

Example:

```go
func TestResourceUser_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:          func() { testAccPreCheck(t) },
        ProviderFactories: providerFactories,
        CheckDestroy:      testAccCheckUserDestroy,
        Steps: []resource.TestStep{
            // Test steps here
        },
    })
}
```

### Skipping Tests

Tests can be skipped in short mode using the `skipIfShortTest` helper function:

```go
func TestResourceUser_basic(t *testing.T) {
    skipIfShortTest(t)
    // Test implementation
}
```

## Adding New Tests

When adding new tests, follow these guidelines:

1. For unit tests, use the standard Go testing package
2. For acceptance tests, use the Terraform SDK's testing framework
3. Always skip acceptance tests when running in short mode
4. Use descriptive test names that indicate what's being tested
5. Use `ProviderFactories` instead of `Providers` in `resource.TestCase`

## Mocking

For unit tests that don't require a running API, you can use the mock HTTP server provided in `client_test.go`. This allows you to test the provider's functionality without making actual API calls.

Example:

```go
// Setup mock server
server := setupMockServer(t, "/api/users/1", http.StatusOK, expectedResponse)
defer server.Close()

// Create client pointing to mock server
client := taskmanager.NewClient(server.URL+"/", "test-token")
```

## Debugging Tests

To debug tests, you can use the `-v` flag for verbose output:

```sh
go test -v ./...
```

For more detailed logging, set the `TF_LOG` environment variable:

```sh
export TF_LOG=DEBUG
go test -v ./...
```

On Windows (PowerShell):

```powershell
$env:TF_LOG="DEBUG"
go test -v ./...
```

## Common Issues and Solutions

### Type Mismatch Errors

If you encounter type mismatch errors related to `Providers` vs `ProviderFactories`, ensure you're using `ProviderFactories` in your `resource.TestCase` structs. This is a requirement for Terraform SDK v2.

### Import Errors

Ensure all necessary packages are imported, especially when using SDK v2 features. Common imports include:

```go
import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)
```

### API Connection Issues

If tests fail due to API connection issues, check:

1. The TaskManager-Go API is running
2. The `TASKMANAGER_BASE_URL` environment variable is set correctly
3. The `TASKMANAGER_TOKEN` environment variable contains a valid token

## Continuous Integration

When setting up CI for this provider, consider:

1. Running unit tests on every PR
2. Running acceptance tests on a schedule or before releases
3. Setting up the TaskManager-Go API as part of the CI pipeline for acceptance tests