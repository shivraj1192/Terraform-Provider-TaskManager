# TaskManager Terraform Provider Documentation

## Table of Contents

- [Introduction](#introduction)
- [Provider Configuration](#provider-configuration)
- [Basic Concepts](#basic-concepts)
- [Resources](#resources)
  - [User Resource](#user-resource)
  - [Team Resource](#team-resource)
  - [Task Resource](#task-resource)
  - [Comment Resource](#comment-resource)
  - [Attachment Resource](#attachment-resource)
- [Data Sources](#data-sources)
  - [User Data Source](#user-data-source)
  - [Team Data Source](#team-data-source)
  - [Task Data Source](#task-data-source)
  - [Comment Data Source](#comment-data-source)
  - [Attachment Data Source](#attachment-data-source)
- [Complete Examples](#complete-examples)
  - [Project Setup Example](#project-setup-example)
  - [Advanced Project Management Example](#advanced-project-management-example)
- [Best Practices](#best-practices)
  - [Security](#security)
  - [Organization](#organization)
  - [State Management](#state-management)
- [Troubleshooting](#troubleshooting)
  - [Common Issues](#common-issues)
  - [Debugging Tips](#debugging-tips)

## Introduction

The TaskManager Terraform Provider allows you to manage your TaskManager resources using infrastructure as code. This provider interacts with the TaskManager-Go API to create, read, update, and delete resources such as users, teams, tasks, comments, and attachments.

### Prerequisites

Before using this provider, ensure you have:

- [Terraform](https://www.terraform.io/downloads.html) v1.0+ installed
- Access to a running [TaskManager-Go API](https://github.com/shivraj1192/TaskManager-Go.git) instance
- An API token for authentication

## Provider Configuration

### Basic Configuration

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
  base_url = "http://localhost:8080/"  # URL of your TaskManager-Go API
  token    = "your-api-token"         # Your API authentication token
}
```

### Provider Arguments

- `base_url` (Required) - The URL of your TaskManager-Go API instance
- `token` (Required) - Your API authentication token

> **Security Note:** Never store your API token directly in your Terraform files. Use environment variables or Terraform variables instead.

## Basic Concepts

### HCL Syntax Basics

HashiCorp Configuration Language (HCL) is used to define resources in Terraform. Here's a quick overview of the syntax:

```hcl
# Block syntax
block_type "label" "name" {
  attribute = value
  nested_block {
    nested_attribute = nested_value
  }
}

# Resource example
resource "taskmanager_user" "example" {
  name  = "John Doe"
  email = "john@example.com"
}
```

### Resource Lifecycle

Each resource in Terraform goes through a lifecycle:

1. **Create** - Resource is created when first applied
2. **Read** - Resource state is read from the API
3. **Update** - Resource is updated when configuration changes
4. **Delete** - Resource is deleted when removed from configuration

### Dependencies

Terraform automatically handles dependencies between resources. For example, if a task depends on a user, Terraform will create the user first.

```hcl
resource "taskmanager_user" "user1" {
  name  = "John Doe"
  email = "john@example.com"
}

resource "taskmanager_task" "task1" {
  title      = "Complete documentation"
  creator_id = taskmanager_user.user1.id  # Reference to user1's ID
}
```

## Resources

### User Resource

The `taskmanager_user` resource allows you to manage users in TaskManager.

#### Example Usage

```hcl
resource "taskmanager_user" "admin" {
  uname    = "admin_user"
  name     = "Admin User"
  email    = "admin@example.com"
  password = "securePassword123"
  role     = "Admin"
}
```

#### Argument Reference

- `uname` (Required) - The username
- `name` (Required) - The full name of the user
- `email` (Required) - The email address of the user
- `password` (Required) - The user's password
- `role` (Required) - The user's role ("Admin" or "Member")

#### Attribute Reference

- `id` - The ID of the user

### Team Resource

The `taskmanager_team` resource allows you to manage teams in TaskManager.

#### Example Usage

```hcl
resource "taskmanager_team" "engineering" {
  name        = "Engineering Team"
  description = "Software engineering team"
  members     = [taskmanager_user.dev1.id, taskmanager_user.dev2.id]
}
```

#### Argument Reference

- `name` (Required) - The name of the team
- `description` (Optional) - A description of the team
- `members` (Optional) - A list of user IDs who are members of this team

#### Attribute Reference

- `id` - The ID of the team

### Task Resource

The `taskmanager_task` resource allows you to manage tasks in TaskManager.

#### Example Usage

```hcl
resource "taskmanager_task" "feature_task" {
  title        = "Implement new feature"
  description  = "Implement the new authentication feature"
  priority     = "High"
  status       = "In progress"
  due_date     = "2023-12-31T23:59:59Z"
  team_id      = taskmanager_team.engineering.id
  creator_id   = taskmanager_user.admin.id
  assignees    = [taskmanager_user.dev1.id, taskmanager_user.dev2.id]
}
```

#### Argument Reference

- `title` (Required) - The title of the task
- `description` (Optional) - A description of the task
- `status` (Optional) - The current status of the task
- `priority` (Optional) - The priority of the task
- `due_date` (Optional) - The due date of the task in RFC3339 format
- `creator_id` (Optional) - The ID of the user who created the task
- `team_id` (Required) - The ID of the team the task belongs to
- `parent_task_id` (Optional) - The ID of the parent task if this is a subtask
- `assignees` (Optional) - A list of user IDs assigned to this task
- `labels` (Optional) - A list of label IDs associated with this task

#### Attribute Reference

- `id` - The ID of the task
- `subtasks` - A list of subtask IDs
- `comments` - A list of comment IDs
- `attachments` - A list of attachment IDs

### Comment Resource

The `taskmanager_comment` resource allows you to manage comments on tasks.

#### Example Usage

```hcl
resource "taskmanager_comment" "status_update" {
  content          = "Making good progress on this task"
  task_id          = taskmanager_task.feature_task.id
  user_id          = taskmanager_user.dev1.id
  parent_comment_id = 0  # 0 means no parent comment
}
```

#### Argument Reference

- `content` (Required) - The content of the comment
- `task_id` (Required) - The ID of the task this comment belongs to
- `user_id` (Required) - The ID of the user who created the comment
- `parent_comment_id` (Optional) - The ID of the parent comment if this is a reply

#### Attribute Reference

- `id` - The ID of the comment

### Attachment Resource

The `taskmanager_attachment` resource allows you to manage file attachments on tasks.

#### Example Usage

```hcl
resource "taskmanager_attachment" "design_doc" {
  file_name   = "design_document.pdf"
  task_id     = taskmanager_task.feature_task.id
  url         = "./static/files/design_document.pdf"
  uploader_id = taskmanager_user.admin.id
}
```

#### Argument Reference

- `file_name` (Required) - The name of the file
- `task_id` (Required) - The ID of the task this attachment belongs to
- `url` (Required) - The URL or path to the file
- `uploader_id` (Required) - The ID of the user who uploaded the attachment

#### Attribute Reference

- `id` - The ID of the attachment

## Data Sources

Data sources allow you to fetch existing resources from the TaskManager API.

### User Data Source

```hcl
data "taskmanager_user" "existing_user" {
  id = 1
}

output "user_email" {
  value = data.taskmanager_user.existing_user.email
}
```

### Team Data Source

```hcl
data "taskmanager_team" "existing_team" {
  id = 1
}

output "team_members" {
  value = data.taskmanager_team.existing_team.members
}
```

### Task Data Source

```hcl
data "taskmanager_task" "existing_task" {
  id = 1
}

output "task_assignees" {
  value = data.taskmanager_task.existing_task.assignees
}
```

### Comment Data Source

```hcl
data "taskmanager_comment" "existing_comment" {
  id = 1
}

output "comment_content" {
  value = data.taskmanager_comment.existing_comment.content
}
```

### Attachment Data Source

```hcl
data "taskmanager_attachment" "existing_attachment" {
  id = 1
}

output "attachment_url" {
  value = data.taskmanager_attachment.existing_attachment.url
}
```

## Complete Examples

### Project Setup Example

This example sets up a complete project structure with teams, users, tasks, comments, and attachments.

```hcl
# Provider configuration
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
  token    = var.api_token
}

# Variables
variable "api_token" {
  description = "API token for TaskManager"
  type        = string
  sensitive   = true
}

# Create users
resource "taskmanager_user" "project_manager" {
  uname    = "pm_user"
  name     = "Project Manager"
  email    = "pm@example.com"
  password = "securePass123"
  role     = "Admin"
}

resource "taskmanager_user" "developer1" {
  uname    = "dev1"
  name     = "Developer One"
  email    = "dev1@example.com"
  password = "securePass456"
  role     = "Member"
}

resource "taskmanager_user" "developer2" {
  uname    = "dev2"
  name     = "Developer Two"
  email    = "dev2@example.com"
  password = "securePass789"
  role     = "Member"
}

# Create a team
resource "taskmanager_team" "project_team" {
  name        = "Project Alpha Team"
  description = "Team responsible for Project Alpha"
  members     = [
    taskmanager_user.project_manager.id,
    taskmanager_user.developer1.id,
    taskmanager_user.developer2.id
  ]
}

# Create parent task
resource "taskmanager_task" "parent_task" {
  title       = "Project Alpha Implementation"
  description = "Complete implementation of Project Alpha"
  priority    = "High"
  status      = "In progress"
  due_date    = "2023-12-31T23:59:59Z"
  team_id     = taskmanager_team.project_team.id
  creator_id  = taskmanager_user.project_manager.id
  assignees   = [
    taskmanager_user.developer1.id,
    taskmanager_user.developer2.id
  ]
}

# Create subtasks
resource "taskmanager_task" "subtask1" {
  title          = "Frontend Implementation"
  description    = "Implement the frontend components"
  priority       = "Medium"
  status         = "To do"
  due_date       = "2023-11-30T23:59:59Z"
  team_id        = taskmanager_team.project_team.id
  creator_id     = taskmanager_user.project_manager.id
  parent_task_id = taskmanager_task.parent_task.id
  assignees      = [taskmanager_user.developer1.id]
}

resource "taskmanager_task" "subtask2" {
  title          = "Backend Implementation"
  description    = "Implement the backend services"
  priority       = "Medium"
  status         = "To do"
  due_date       = "2023-11-30T23:59:59Z"
  team_id        = taskmanager_team.project_team.id
  creator_id     = taskmanager_user.project_manager.id
  parent_task_id = taskmanager_task.parent_task.id
  assignees      = [taskmanager_user.developer2.id]
}

# Add comments
resource "taskmanager_comment" "kickoff_comment" {
  content   = "Let's start working on this project. Please check the attached requirements document."
  task_id   = taskmanager_task.parent_task.id
  user_id   = taskmanager_user.project_manager.id
}

# Add attachments
resource "taskmanager_attachment" "requirements_doc" {
  file_name   = "requirements.pdf"
  task_id     = taskmanager_task.parent_task.id
  url         = "./static/files/requirements.pdf"
  uploader_id = taskmanager_user.project_manager.id
}

# Outputs
output "project_team_id" {
  value = taskmanager_team.project_team.id
}

output "parent_task_id" {
  value = taskmanager_task.parent_task.id
}
```

### Advanced Project Management Example

This comprehensive example demonstrates all the concepts covered in this documentation, including:
- Provider configuration with variables
- Resource creation and dependencies
- Data sources for existing resources
- Modules for organization
- State management
- Output values

```hcl
# main.tf - Main configuration file

terraform {
  required_providers {
    taskmanager = {
      source  = "local/taskmanager/taskmanager"
      version = "0.1.0"
    }
  }
  
  # Remote state configuration for team collaboration
  backend "s3" {
    bucket         = "taskmanager-terraform-state"
    key            = "project-management/terraform.tfstate"
    region         = "us-west-2"
    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}

provider "taskmanager" {
  base_url = var.taskmanager_api_url
  token    = var.taskmanager_api_token
}

# Import local variable definitions
locals {
  environment = "production"
  project_name = "Enterprise Task Management System"
  common_tags = {
    Project     = local.project_name
    Environment = local.environment
    ManagedBy   = "Terraform"
  }
}

# Create departments using the team module
module "engineering_department" {
  source = "./modules/team"
  
  team_name        = "Engineering"
  team_description = "Software Engineering Department"
  admin_users      = [module.users.admin_user_ids]
  member_users     = [module.users.developer_user_ids]
  tags             = merge(local.common_tags, { Department = "Engineering" })
}

module "product_department" {
  source = "./modules/team"
  
  team_name        = "Product"
  team_description = "Product Management Department"
  admin_users      = [module.users.product_manager_ids]
  member_users     = [module.users.designer_user_ids]
  tags             = merge(local.common_tags, { Department = "Product" })
}

# Create users using the users module
module "users" {
  source = "./modules/users"
  
  admin_users = [
    {
      username = "admin1"
      name     = "Admin User"
      email    = "admin@example.com"
      role     = "Admin"
    }
  ]
  
  developer_users = [
    {
      username = "dev1"
      name     = "Developer One"
      email    = "dev1@example.com"
      role     = "Member"
    },
    {
      username = "dev2"
      name     = "Developer Two"
      email    = "dev2@example.com"
      role     = "Member"
    }
  ]
  
  product_manager_users = [
    {
      username = "pm1"
      name     = "Product Manager"
      email    = "pm@example.com"
      role     = "Admin"
    }
  ]
  
  designer_users = [
    {
      username = "designer1"
      name     = "UI Designer"
      email    = "designer@example.com"
      role     = "Member"
    }
  ]
  
  password_policy = var.password_policy
  tags            = local.common_tags
}

# Create project tasks
module "project_tasks" {
  source = "./modules/tasks"
  
  project_name    = "Mobile App Development"
  team_id         = module.engineering_department.team_id
  creator_id      = module.users.admin_user_ids[0]
  assignee_ids    = module.users.developer_user_ids
  start_date      = "2023-01-01T00:00:00Z"
  end_date        = "2023-12-31T23:59:59Z"
  priority_levels = ["High", "Medium", "Low"]
  tags            = local.common_tags
}

# Query existing resources using data sources
data "taskmanager_team" "existing_team" {
  id = 1  # ID of an existing team
}

data "taskmanager_user" "existing_user" {
  id = 1  # ID of an existing user
}

# Create a task assigned to an existing user
resource "taskmanager_task" "existing_user_task" {
  title       = "Task for existing user"
  description = "This task is assigned to an existing user"
  priority    = "Medium"
  status      = "To do"
  due_date    = timeadd(timestamp(), "168h")  # 1 week from now
  team_id     = data.taskmanager_team.existing_team.id
  creator_id  = module.users.admin_user_ids[0]
  assignees   = [data.taskmanager_user.existing_user.id]
}

# Conditional resource creation
resource "taskmanager_task" "conditional_task" {
  count = var.create_extra_tasks ? 1 : 0
  
  title       = "Conditional Task"
  description = "This task is created conditionally"
  priority    = "Low"
  status      = "To do"
  team_id     = module.engineering_department.team_id
  creator_id  = module.users.admin_user_ids[0]
}

# Dynamic block example for multiple attachments
resource "taskmanager_task" "task_with_attachments" {
  title       = "Task with multiple attachments"
  description = "This task has multiple attachments added dynamically"
  priority    = "Medium"
  status      = "To do"
  team_id     = module.engineering_department.team_id
  creator_id  = module.users.admin_user_ids[0]
}

# Dynamic attachment creation
resource "taskmanager_attachment" "dynamic_attachments" {
  for_each = toset(var.attachment_files)
  
  file_name   = each.key
  task_id     = taskmanager_task.task_with_attachments.id
  url         = "./static/files/${each.key}"
  uploader_id = module.users.admin_user_ids[0]
}

# Output important information
output "engineering_team_id" {
  description = "ID of the Engineering team"
  value       = module.engineering_department.team_id
}

output "product_team_id" {
  description = "ID of the Product team"
  value       = module.product_department.team_id
}

output "admin_users" {
  description = "List of admin users"
  value       = module.users.admin_user_details
  sensitive   = true  # Mask sensitive information
}

output "project_task_ids" {
  description = "IDs of all project tasks"
  value       = module.project_tasks.task_ids
}
```

```hcl
# variables.tf - Variable definitions

variable "taskmanager_api_url" {
  description = "URL of the TaskManager API"
  type        = string
  default     = "http://localhost:8080/"
}

variable "taskmanager_api_token" {
  description = "Authentication token for the TaskManager API"
  type        = string
  sensitive   = true
}

variable "password_policy" {
  description = "Password policy configuration"
  type = object({
    min_length      = number
    require_numbers = bool
    require_symbols = bool
  })
  default = {
    min_length      = 8
    require_numbers = true
    require_symbols = true
  }
}

variable "create_extra_tasks" {
  description = "Whether to create additional tasks"
  type        = bool
  default     = false
}

variable "attachment_files" {
  description = "List of attachment files to add to tasks"
  type        = list(string)
  default     = ["requirements.pdf", "design.pdf", "schedule.xlsx"]
}
```

```hcl
# modules/team/main.tf - Team module

variable "team_name" {
  description = "Name of the team"
  type        = string
}

variable "team_description" {
  description = "Description of the team"
  type        = string
  default     = ""
}

variable "admin_users" {
  description = "List of admin user IDs"
  type        = list(number)
  default     = []
}

variable "member_users" {
  description = "List of member user IDs"
  type        = list(number)
  default     = []
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

resource "taskmanager_team" "team" {
  name        = var.team_name
  description = var.team_description
  members     = concat(var.admin_users, var.member_users)
}

output "team_id" {
  description = "ID of the created team"
  value       = taskmanager_team.team.id
}

output "team_name" {
  description = "Name of the created team"
  value       = taskmanager_team.team.name
}
```

```hcl
# modules/users/main.tf - Users module

variable "admin_users" {
  description = "List of admin users to create"
  type = list(object({
    username = string
    name     = string
    email    = string
    role     = string
  }))
  default = []
}

variable "developer_users" {
  description = "List of developer users to create"
  type = list(object({
    username = string
    name     = string
    email    = string
    role     = string
  }))
  default = []
}

variable "product_manager_users" {
  description = "List of product manager users to create"
  type = list(object({
    username = string
    name     = string
    email    = string
    role     = string
  }))
  default = []
}

variable "designer_users" {
  description = "List of designer users to create"
  type = list(object({
    username = string
    name     = string
    email    = string
    role     = string
  }))
  default = []
}

variable "password_policy" {
  description = "Password policy configuration"
  type = object({
    min_length      = number
    require_numbers = bool
    require_symbols = bool
  })
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

# Generate a secure password based on policy
locals {
  password_suffix = var.password_policy.require_symbols ? "!2#4" : "1234"
}

# Create admin users
resource "taskmanager_user" "admin_users" {
  count = length(var.admin_users)
  
  uname    = var.admin_users[count.index].username
  name     = var.admin_users[count.index].name
  email    = var.admin_users[count.index].email
  password = "${var.admin_users[count.index].username}${local.password_suffix}"
  role     = var.admin_users[count.index].role
}

# Create developer users
resource "taskmanager_user" "developer_users" {
  count = length(var.developer_users)
  
  uname    = var.developer_users[count.index].username
  name     = var.developer_users[count.index].name
  email    = var.developer_users[count.index].email
  password = "${var.developer_users[count.index].username}${local.password_suffix}"
  role     = var.developer_users[count.index].role
}

# Create product manager users
resource "taskmanager_user" "product_manager_users" {
  count = length(var.product_manager_users)
  
  uname    = var.product_manager_users[count.index].username
  name     = var.product_manager_users[count.index].name
  email    = var.product_manager_users[count.index].email
  password = "${var.product_manager_users[count.index].username}${local.password_suffix}"
  role     = var.product_manager_users[count.index].role
}

# Create designer users
resource "taskmanager_user" "designer_users" {
  count = length(var.designer_users)
  
  uname    = var.designer_users[count.index].username
  name     = var.designer_users[count.index].name
  email    = var.designer_users[count.index].email
  password = "${var.designer_users[count.index].username}${local.password_suffix}"
  role     = var.designer_users[count.index].role
}

# Outputs
output "admin_user_ids" {
  description = "IDs of admin users"
  value       = taskmanager_user.admin_users[*].id
}

output "developer_user_ids" {
  description = "IDs of developer users"
  value       = taskmanager_user.developer_users[*].id
}

output "product_manager_ids" {
  description = "IDs of product manager users"
  value       = taskmanager_user.product_manager_users[*].id
}

output "designer_user_ids" {
  description = "IDs of designer users"
  value       = taskmanager_user.designer_users[*].id
}

output "admin_user_details" {
  description = "Details of admin users"
  value = [
    for user in taskmanager_user.admin_users : {
      id    = user.id
      name  = user.name
      email = user.email
      role  = user.role
    }
  ]
}
```

```hcl
# modules/tasks/main.tf - Tasks module

variable "project_name" {
  description = "Name of the project"
  type        = string
}

variable "team_id" {
  description = "ID of the team"
  type        = number
}

variable "creator_id" {
  description = "ID of the task creator"
  type        = number
}

variable "assignee_ids" {
  description = "List of assignee IDs"
  type        = list(number)
  default     = []
}

variable "start_date" {
  description = "Project start date"
  type        = string
}

variable "end_date" {
  description = "Project end date"
  type        = string
}

variable "priority_levels" {
  description = "List of priority levels"
  type        = list(string)
  default     = ["High", "Medium", "Low"]
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

# Create main project task
resource "taskmanager_task" "main_task" {
  title       = var.project_name
  description = "Main task for ${var.project_name}"
  priority    = "High"
  status      = "In progress"
  due_date    = var.end_date
  team_id     = var.team_id
  creator_id  = var.creator_id
  assignees   = var.assignee_ids
}

# Create planning phase task
resource "taskmanager_task" "planning_task" {
  title          = "Planning Phase"
  description    = "Planning phase for ${var.project_name}"
  priority       = var.priority_levels[0]  # High
  status         = "In progress"
  due_date       = timeadd(var.start_date, "336h")  # 2 weeks after start
  team_id        = var.team_id
  creator_id     = var.creator_id
  parent_task_id = taskmanager_task.main_task.id
  assignees      = [var.assignee_ids[0]]
}

# Create development phase task
resource "taskmanager_task" "development_task" {
  title          = "Development Phase"
  description    = "Development phase for ${var.project_name}"
  priority       = var.priority_levels[0]  # High
  status         = "To do"
  due_date       = timeadd(var.start_date, "1344h")  # 8 weeks after start
  team_id        = var.team_id
  creator_id     = var.creator_id
  parent_task_id = taskmanager_task.main_task.id
  assignees      = var.assignee_ids
}

# Create testing phase task
resource "taskmanager_task" "testing_task" {
  title          = "Testing Phase"
  description    = "Testing phase for ${var.project_name}"
  priority       = var.priority_levels[1]  # Medium
  status         = "To do"
  due_date       = timeadd(var.start_date, "1680h")  # 10 weeks after start
  team_id        = var.team_id
  creator_id     = var.creator_id
  parent_task_id = taskmanager_task.main_task.id
  assignees      = var.assignee_ids
}

# Create deployment phase task
resource "taskmanager_task" "deployment_task" {
  title          = "Deployment Phase"
  description    = "Deployment phase for ${var.project_name}"
  priority       = var.priority_levels[0]  # High
  status         = "To do"
  due_date       = timeadd(var.start_date, "2016h")  # 12 weeks after start
  team_id        = var.team_id
  creator_id     = var.creator_id
  parent_task_id = taskmanager_task.main_task.id
  assignees      = [var.assignee_ids[0]]
}

# Outputs
output "task_ids" {
  description = "IDs of all created tasks"
  value = {
    main_task       = taskmanager_task.main_task.id
    planning_task   = taskmanager_task.planning_task.id
    development_task = taskmanager_task.development_task.id
    testing_task    = taskmanager_task.testing_task.id
    deployment_task = taskmanager_task.deployment_task.id
  }
}

## Best Practices

### Security

1. **Never hardcode API tokens** in your Terraform files. Use environment variables or Terraform variables instead.

   ```hcl
   provider "taskmanager" {
     base_url = "http://localhost:8080/"
     token    = var.api_token  # Define this in a variables.tf file
   }
   ```

2. **Use sensitive = true** for variables containing sensitive information.

   ```hcl
   variable "api_token" {
     description = "API token for TaskManager"
     type        = string
     sensitive   = true  # Masks the value in logs and outputs
   }
   ```

### Organization

1. **Use modules** for reusable components.

   ```hcl
   module "project_team" {
     source = "./modules/team"
     
     team_name = "Engineering"
     members   = [1, 2, 3]
   }
   ```

2. **Use consistent naming conventions** for resources.

   ```hcl
   resource "taskmanager_user" "dev_john" {...}
   resource "taskmanager_user" "dev_jane" {...}
   resource "taskmanager_team" "team_engineering" {...}
   ```

3. **Group related resources** together in your configuration files.

### State Management

1. **Use remote state** for team collaboration.

   ```hcl
   terraform {
     backend "s3" {
       bucket = "my-terraform-state"
       key    = "taskmanager/terraform.tfstate"
       region = "us-west-2"
     }
   }
   ```

2. **Lock your state file** to prevent concurrent modifications.

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   - Ensure your API token is valid and not expired
   - Check that the base_url is correct and the API is running

2. **Resource Creation Failures**
   - Check the API logs for detailed error messages
   - Ensure all required fields are provided
   - Verify that referenced resources exist

3. **State Drift**
   - If resources were modified outside of Terraform, run `terraform refresh` to update the state
   - Use `terraform import` to bring existing resources under Terraform management

### Debugging Tips

1. **Enable Terraform Logging**

   ```sh
   export TF_LOG=DEBUG
   export TF_LOG_PATH=./terraform.log
   ```

2. **Use terraform plan** to preview changes before applying them.

3. **Check the TaskManager API logs** for detailed error messages.

---

This documentation provides a comprehensive guide to using the TaskManager Terraform Provider. For additional help or to report issues, please refer to the project repository.