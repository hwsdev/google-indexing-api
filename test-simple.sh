#!/bin/bash

# Simple test script for Google Indexing API
BASE_URL="http://localhost:8080"

echo "🚀 Testing Google Indexing API"
echo "================================"

# Test 1: Health Check
echo "1. Testing Health Check..."
curl -s "$BASE_URL/api/health" | jq .
echo ""

# Test 2: Submit Single URL (you need to replace with your service account)
echo "2. Testing Single URL Submission..."
echo "⚠️  Please replace with your actual service account credentials!"

cat << 'EOF'
curl -X POST http://localhost:8080/api/v1/index \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/page",
    "service_account": {
      "type": "service_account",
      "project_id": "your-project-id",
      "private_key": "-----BEGIN PRIVATE KEY-----\nYOUR_KEY_HERE\n-----END PRIVATE KEY-----\n",
      "client_email": "service@project.iam.gserviceaccount.com",
      "client_id": "123456789",
      "auth_uri": "https://accounts.google.com/o/oauth2/auth",
      "token_uri": "https://oauth2.googleapis.com/token",
      "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
      "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/service%40project.iam.gserviceaccount.com"
    }
  }'
EOF

echo ""
echo ""

# Test 3: Cache Stats
echo "3. Testing Cache Stats..."
curl -s "$BASE_URL/api/v1/cache/stats" | jq .
echo ""

echo "✅ Test completed!"
echo "💡 Replace service account credentials in the curl command above to test actual indexing."
