# Terraform Provider: TaskManager

A custom [Terraform](https://www.terraform.io/) provider for managing users, teams, tasks, comments, and attachments in your [TaskManager-Go](https://github.com/shivraj1192/TaskManager-Go.git) backend.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Step-by-Step Manual](#step-by-step-manual)
  - [1. Start the TaskManager-Go API](#1-start-the-taskmanager-go-api)
  - [2. Get Your API Token](#2-get-your-api-token-using-postman)
  - [3. Clone and Build This Provider](#3-clone-and-build-this-provider)
  - [4. Install the Provider Binary](#4-install-the-provider-binary)
  - [5. Test or Create Your Terraform Script](#5-test-or-create-your-terraform-script)
- [Example main.tf Explained](#example-maintf-explained)
- [Resources](#resources)
- [Data Sources](#data-sources)
- [Testing](#testing)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

---

## Overview

This Terraform provider lets you manage users, teams, tasks, comments, and attachments in your TaskManager-Go backend using Terraform.

> **Important:**  
> You must have the [TaskManager-Go API](https://github.com/shivraj1192/TaskManager-Go.git) running before using this provider.  
> See [TaskManager-Go Setup Guide](https://github.com/shivraj1192/TaskManager-Go#readme) for instructions.

---

## Features

- **User Management:** Create, update, delete, and import users.
- **Team Management:** Manage teams, members, and related tasks.
- **Task Management:** Create and assign tasks, set priorities, and manage relationships.
- **Comments & Attachments:** Add comments and upload attachments to tasks.
- **Data Sources:** Query users, teams, tasks, comments, and attachments.

---

## Project Structure

```
.
├── taskmanager/
│   ├── provider.go
│   ├── resource_team.go
│   ├── resource_user.go
│   ├── resource_task.go
│   ├── resource_comment.go
│   ├── resource_attachment.go
│   ├── data_source_team.go
│   ├── data_source_user.go
│   ├── data_source_task.go
│   ├── data_source_comment.go
│   └── data_source_attachment.go
├── test_taskmanager/
│   ├── provider_test.go
│   ├── resource_user_test.go
│   ├── resource_team_test.go
│   ├── resource_task_test.go
│   ├── resource_comment_test.go
│   ├── resource_attachment_test.go
│   ├── data_source_test.go
│   ├── test_helper.go
│   └── Provider_Test.md
├── test_taskmanager_with_HCL/
│   ├── static
│   │   └── files
│   └── main.tf
├── main.go
├── go.mod
├── go.sum
├── README.md
└── Terraform_Provider_Taskmanager_Doc.md
```

---

## Prerequisites

- [Go](https://go.dev/doc/install) 1.18 or later
- [Terraform](https://www.terraform.io/downloads.html) v1.0+
- Access to a running [TaskManager-Go API](https://github.com/shivraj1192/TaskManager-Go.git) instance

---

## Step-by-Step Manual

### 1. Start the TaskManager-Go API

You **must** have the TaskManager-Go backend running before using this provider.

- Follow the setup guide here:  
  👉 [TaskManager-Go Setup Guide](https://github.com/shivraj1192/TaskManager-Go#readme)

#### Example (Windows, Linux, Mac):

```sh
# Clone the TaskManager-Go project
git clone https://github.com/shivraj1192/TaskManager-Go.git
cd TaskManager-Go

# Install dependencies
go mod tidy

# Start the API server
go run ./cmd/main.go
# or use 'air' if available for live reload
```

---

### 2. Get Your API Token (using Postman)

Before using the provider, you need an API token from the TaskManager-Go backend.  
Follow these steps using [Postman](https://www.postman.com/) or any API client:

**a) Register a User**

- **POST** `http://localhost:8080/api/register`
- **Body (JSON):**
  ```json
  {
    "uname": "yourusername",
    "name": "Your Name",
    "email": "your@email.com",
    "password": "yourpassword"
  }
  ```

**b) Login to Get Token**

- **POST** `http://localhost:8080/api/login`
- **Body (JSON):**
  ```json
  {
    "email": "your@email.com",
    "password": "yourpassword"
  }
  ```
- **Response:**  
  The response will include a field like `"token": "eyJhbGciOi..."`.  
  **Copy this token.**

**c) Use the Token in Terraform**

Paste the token in your `main.tf` provider block:

```hcl
provider "taskmanager" {
  base_url = "http://localhost:8080/"
  token    = "PASTE_YOUR_TOKEN_HERE"
}
```

---

### 3. Clone and Build This Provider

Open a new terminal and clone this provider project:

```sh
git clone https://github.com/your-org/terraform-provider-taskmanager.git
cd terraform-provider-taskmanager
go mod tidy
```

#### Build the provider binary:

- **Windows:**
  ```sh
  go build -o terraform-provider-taskmanager.exe
  ```
- **Linux/Mac:**
  ```sh
  go build -o terraform-provider-taskmanager
  ```

---

### 4. Install the Provider Binary

Move the binary to your Terraform plugins directory:

- **Windows (PowerShell):**
  ```powershell
  $dest = "$env:APPDATA\terraform.d\plugins\local\taskmanager\taskmanager\0.1.0\windows_amd64"
  New-Item -ItemType Directory -Force -Path $dest
  Move-Item -Path ".\terraform-provider-taskmanager.exe" -Destination "$dest\terraform-provider-taskmanager.exe"
  ```
- **Linux:**
  ```sh
  mkdir -p ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/linux_amd64
  mv terraform-provider-taskmanager ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/linux_amd64/
  ```
- **Mac:**
  ```sh
  mkdir -p ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_amd64
  mv terraform-provider-taskmanager ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_amd64/
  ```

#### To delete the plugin (cleanup):

- **Windows (PowerShell):**
  ```powershell
  Remove-Item -Path "$env:APPDATA\terraform.d\plugins\local\taskmanager\taskmanager\0.1.0\windows_amd64" -Recurse -Force
  ```
- **Linux:**
  ```sh
  rm -rf ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/linux_amd64
  ```
- **Mac:**
  ```sh
  rm -rf ~/.terraform.d/plugins/local/taskmanager/taskmanager/0.1.0/darwin_amd64
  ```

---

### 5. Test or Create Your Terraform Script

Navigate to the `test_taskmanager_with_HCL` directory (provided in this repo):

```sh
cd test_taskmanager_with_HCL
```

You can **use the provided `main.tf` file** directly, or **create your own Terraform script**.  
If you want to create your own, see [`Terraform_Provider_Taskmanager_Doc.md`](Terraform_Provider_Taskmanager_Doc.md) for detailed documentation and examples.

To initialize and apply:

```sh
terraform init
terraform plan
terraform apply
```

---

## Example `main.tf` Explained

Below is a sample `main.tf` file for using the TaskManager Terraform provider:

```hcl
terraform {
  required_providers {
    taskmanager = {
      source  = "local/taskmanager/taskmanager"
      version = "0.1.0"
    }
  }
}

provider "taskmanager" {
  base_url = "http://localhost:8080/"
  token    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NDk2MTgwNzUsInVzZXJfaWQiOjF9.ZctJvp90AVQFYU55TCfhtGscWrHPMoyk-GCamtCOa9A"
}
```

#### Explanation

- **terraform block**:  
  Specifies the required provider (`taskmanager`) and its source/version.  
  This tells Terraform to use your local TaskManager provider plugin.

- **provider "taskmanager" block**:  
  - `base_url`: The URL where your TaskManager-Go API is running (default: `http://localhost:8080/`).
  - `token`: The API token for authenticating requests.  
    You must generate this token by registering and logging in to your TaskManager-Go API (see the guide above for how to get it using Postman).

**Note:**  
Never commit your real API token to a public repository.  
For more details on available resources and configuration, see the [Terraform_Provider_Taskmanager_Doc.md](Terraform_Provider_Taskmanager_Doc.md).

---

## Resources

The TaskManager provider includes the following resources:

- `taskmanager_user`: Manage users in the TaskManager system
- `taskmanager_team`: Create and manage teams and their members
- `taskmanager_task`: Create, assign, and manage tasks
- `taskmanager_comment`: Add comments to tasks
- `taskmanager_attachment`: Upload and manage file attachments for tasks

For detailed documentation on each resource, see the [Terraform_Provider_Taskmanager_Doc.md](Terraform_Provider_Taskmanager_Doc.md).

## Data Sources

The TaskManager provider includes the following data sources:

- `taskmanager_user`: Query user information
- `taskmanager_team`: Query team information and members
- `taskmanager_task`: Query task details and assignments
- `taskmanager_comment`: Query comments on tasks
- `taskmanager_attachment`: Query file attachments for tasks

For detailed documentation on each data source, see the [Terraform_Provider_Taskmanager_Doc.md](Terraform_Provider_Taskmanager_Doc.md).

---

## Testing

This provider includes a comprehensive test suite to ensure functionality and compatibility with the TaskManager-Go API.

### Running Tests

```sh
# Run all tests
go test -v ./...

# Run only unit tests (no API required)
go test -v -short ./...

# Run specific tests
go test -v ./test_taskmanager -run TestProvider
```

### Test Documentation

For detailed information about testing the provider, including test structure, environment setup, and troubleshooting, see the [Provider_Test.md](test_taskmanager/Provider_Test.md) documentation.

---

## Development

1. Clone the repo and install dependencies:

   ```sh
   git clone https://github.com/your-org/terraform-provider-taskmanager.git
   cd terraform-provider-taskmanager
   go mod tidy
   ```

2. Run tests:

   ```sh
   go test ./...
   ```

3. Make changes to the provider code in the `taskmanager/` directory.

4. Build and install the provider as described in the [Step-by-Step Manual](#step-by-step-manual).

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests to ensure they pass
5. Submit a pull request

---

## License

This project is licensed under the MIT License.

---

## Acknowledgements

- [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk)
- [TaskManager-Go](https://github.com/shivraj1192/TaskManager-Go.git)

---

_This project is not affiliated with HashiCorp or Terraform._