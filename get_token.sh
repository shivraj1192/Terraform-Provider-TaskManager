#!/bin/bash

API_URL="http://localhost:8080/api"

# Register test user
curl -s -X POST "$API_URL/register" -H "Content-Type: application/json" -d '{
  "uname": "AT-Test",
  "name": "Test",
  "email": "Test@example.com",
  "password": "TestPass@1234",
  "role": "Admin"
}' > /dev/null

# Login and extract token
TOKEN=$(curl -s -X POST "$API_URL/login" -H "Content-Type: application/json" -d '{
  "uname": "AT-Test",
  "password": "TestPass@1234"
}' | jq -r '.token')

# Export for GitHub Actions
echo "TOKEN=$TOKEN" >> $GITHUB_ENV
