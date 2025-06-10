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

resource "taskmanager_team" "team1" {
  name        = "Team One1"
  description = "Test team"
  members     = [taskmanager_user.adminuser1.id,taskmanager_user.user1.id]
}

resource "taskmanager_user" "user1" {
  uname = "a"
  name = "a"
  email = "a@gmail.com"
  password = "a@1192"
  role = "Member"
}

resource "taskmanager_user" "adminuser1" {
  uname = "c"
  name = "c"
  email = "c@example.com"
  password = "cPass@1234"
  role = "Admin"
}

resource "taskmanager_task" "task1" {
  title = "a"
  description = "a"
  priority = "High"
  status = "In progress"
  due_date = "2025-06-25T18:30:00Z"
  team_id = taskmanager_team.team1.id
  creator_id = taskmanager_user.adminuser1.id
  parent_task_id = 0
  assignees = [taskmanager_user.user1.id,taskmanager_user.user1.id]
  labels = []
}

resource "taskmanager_comment" "comment1"{
  content = "comment1 content"
  task_id = taskmanager_task.task1.id
  user_id = taskmanager_user.adminuser1.id
  parent_comment_id = 0
}

resource "taskmanager_attachment" "attachment2"{
  file_name = "infralovers_courses.pdf"
  task_id = taskmanager_task.task1.id
  url = "./static/files/infralovers_courses.pdf"
  uploader_id = taskmanager_user.adminuser1.id
}

data "taskmanager_user" "user_data1" {
  id = 4
}

data "taskmanager_team" "team_data1" {
  id = 1
}

data "taskmanager_task" "task_data1"{
  id = 1
}

data "taskmanager_comment" "comment_data1" {
  id = 1
}

data "taskmanager_attachment" "attachment_data1" {
  id = 1
}