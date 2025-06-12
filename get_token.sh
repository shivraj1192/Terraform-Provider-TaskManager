#!/bin/bash

API_URL="http://localhost:8080/api"

# Register the test user (ignore if already exists)
curl -s -X POST "$API_URL/register" -H "Content-Type: application/json" -d '{
  "name": "Test User",
  "email": "test@taskmanager.com",
  "password": "test123"
}' > /dev/null

# Login and extract the token
TOKEN=$(curl -s -X POST "$API_URL/login" -H "Content-Type: application/json" -d '{
  "email": "test@taskmanager.com",
  "password": "test123"
}' | jq -r '.token')

# Export token to GitHub Actions environment
echo "TOKEN=$TOKEN" >> $GITHUB_ENV
