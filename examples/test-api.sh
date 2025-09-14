#!/bin/bash

# Google Indexing API Test Scripts
# Make sure to set your API_KEY and BASE_URL before running

BASE_URL="http://localhost:8080"
API_KEY="your-api-key-here"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_test() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Test Health Check
test_health() {
    print_test "Testing Health Check Endpoint"
    
    response=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/api/health")
    
    if [ "$response" -eq 200 ]; then
        print_success "Health check passed"
        curl -s "$BASE_URL/api/health" | jq .
    else
        print_error "Health check failed (HTTP $response)"
    fi
    echo ""
}

# Test Single URL Submission
test_single_url() {
    print_test "Testing Single URL Submission"
    
    response=$(curl -s -X POST "$BASE_URL/api/v1/index" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_KEY" \
        -d '{
            "url": "https://example.com/test-page"
        }')
    
    echo "$response" | jq .
    echo ""
}

# Test Batch URL Submission
test_batch_urls() {
    print_test "Testing Batch URL Submission"
    
    response=$(curl -s -X POST "$BASE_URL/api/v1/index/batch" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_KEY" \
        -d '{
            "urls": [
                "https://example.com/page1",
                "https://example.com/page2",
                "https://example.com/page3"
            ]
        }')
    
    echo "$response" | jq .
    echo ""
}

# Test URL Status Check
test_url_status() {
    print_test "Testing URL Status Check"
    
    encoded_url=$(echo "https://example.com/test-page" | jq -sRr @uri)
    
    response=$(curl -s "$BASE_URL/api/v1/status/$encoded_url" \
        -H "Authorization: Bearer $API_KEY")
    
    echo "$response" | jq .
    echo ""
}

# Test Authentication Error
test_auth_error() {
    print_test "Testing Authentication Error"
    
    response=$(curl -s -w "%{http_code}" "$BASE_URL/api/v1/index" \
        -H "Content-Type: application/json" \
        -d '{
            "url": "https://example.com/test-page"
        }')
    
    echo "Response: $response"
    echo ""
}

# Test Invalid URL
test_invalid_url() {
    print_test "Testing Invalid URL"
    
    response=$(curl -s -X POST "$BASE_URL/api/v1/index" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_KEY" \
        -d '{
            "url": "invalid-url"
        }')
    
    echo "$response" | jq .
    echo ""
}

# Main test runner
run_all_tests() {
    echo "Running Google Indexing API Tests"
    echo "=================================="
    echo ""
    
    test_health
    test_auth_error
    test_invalid_url
    
    if [ "$API_KEY" = "your-api-key-here" ]; then
        print_error "Please set your API_KEY before running authenticated tests"
        exit 1
    fi
    
    test_single_url
    test_batch_urls
    test_url_status
    
    print_success "All tests completed!"
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    print_error "jq is not installed. Please install jq for JSON formatting."
    exit 1
fi

# Check if curl is installed
if ! command -v curl &> /dev/null; then
    print_error "curl is not installed. Please install curl."
    exit 1
fi

# Run tests based on argument
case "${1:-all}" in
    health)
        test_health
        ;;
    single)
        test_single_url
        ;;
    batch)
        test_batch_urls
        ;;
    status)
        test_url_status
        ;;
    auth)
        test_auth_error
        ;;
    invalid)
        test_invalid_url
        ;;
    all)
        run_all_tests
        ;;
    *)
        echo "Usage: $0 {health|single|batch|status|auth|invalid|all}"
        echo ""
        echo "Commands:"
        echo "  health  - Test health check endpoint"
        echo "  single  - Test single URL submission"
        echo "  batch   - Test batch URL submission"
        echo "  status  - Test URL status check"
        echo "  auth    - Test authentication error"
        echo "  invalid - Test invalid URL handling"
        echo "  all     - Run all tests"
        exit 1
        ;;
esac
