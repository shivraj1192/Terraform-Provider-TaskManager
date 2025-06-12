#!/bin/bash

set -e

API_URL="http://localhost:8080/api"
EMAIL="test@taskmanager.com"
PASSWORD="test123"

# Retry logic
retry() {
  local n=0
  local try=$1
  local cmd="${@:2}"
  until [ $n -ge $try ]
  do
    $cmd && break
    n=$((n+1))
    echo "Retrying ($n/$try)..."
    sleep 2
  done
  if [ $n -eq $try ]; then
    echo "Failed after $try attempts"
    exit 1
  fi
}

echo "Registering user..."
retry 5 curl -s -o /dev/null -w "%{http_code}" -X POST "$API_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Test User\", \"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}"

echo "Logging in..."
TOKEN=$(retry 5 curl -s -X POST "$API_URL/login" -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Login failed: token is empty"
  exit 1
fi

echo "✅ Token retrieved"
echo "TOKEN=$TOKEN" >> $GITHUB_ENV
