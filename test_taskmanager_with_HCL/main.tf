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
  token    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NTAxMzk3NjcsInVzZXJfaWQiOjZ9.wUTN0eAmU-hETBO3xZNlaV9pfSLH-Oz8bHUM0JoNNQw"
}

resource "taskmanager_user" "user_new" {
  uname    = "AT - TASKMANAGER UNAME"
  name     = "AT - TASKMANAGER NAME"
  email    = "AT.TASKMANAGER@gmail.com"
  password = "AT - TASKMANAGER PASSWORD"
  role     = "Member"
}

resource "taskmanager_team" "team_new" {
  name        = "AT - TASKMANAGER NAME"
  description = "AT - TASKMANAGER DESC"
  members     = [1, taskmanager_user.user_new.id]
}

resource "taskmanager_task" "task_new" {
  title = "AT - TASKMANAGER TITLE"
  description = "AT - TASKMANAGER DESC"
  priority = "High"
  status = "In progress"
  due_date = "2025-06-25T18:30:00Z"
  team_id = taskmanager_team.team_new.id
  parent_task_id = 0
  assignees = [taskmanager_user.user_new.id]
  labels = []
}

resource "taskmanager_comment" "comment_new"{
  content = "AT - TASKMANAGER CONTENT"
  task_id = taskmanager_task.task_new.id
  parent_comment_id = 0
}

resource "taskmanager_attachment" "attachment_new"{
  file_name = "infralovers_courses.pdf"
  task_id = taskmanager_task.task_new.id
  url = "./static/files/infralovers_courses.pdf"
}