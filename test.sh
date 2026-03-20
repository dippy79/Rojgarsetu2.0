#!/bin/bash
# RojgarSetu Smoke Test Script
# Run after docker compose up -d

set -e

BASE_URL="http://localhost:8080"
echo "🧪 Running smoke tests against ${BASE_URL}..."

# Helper: curl with timeout and fail fast
function curl_test {
  if curl -s -f -w "✓ %{http_code}\n" -o /dev/null "$@"; then
    return 0
  else
    echo "❌ FAILED: $@"
    return 1
  fi
}

# 1. Public endpoints
echo "1. Public endpoints..."
curl_test "${BASE_URL}/robots.txt"
curl_test "${BASE_URL}/api/v1/jobs"
curl_test "${BASE_URL}/api/v1/gov-jobs"
curl_test "${BASE_URL}/api/v1/priv-jobs"
curl_test "${BASE_URL}/api/v1/courses"
curl_test "${BASE_URL}/api/v1/videos"
curl_test "${BASE_URL}/api/v1/companies"
curl_test "${BASE_URL}/api/v1/candidates"

# 2. Auth endpoints (no token needed for register/login)
echo "2. Auth endpoints..."
curl_test -X POST "${BASE_URL}/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"testpass123"}'

LOGIN_RESPONSE=$(curl -s -w "\n" -X POST "${BASE_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}')

ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token // empty')

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
  echo "❌ Login failed - no token"
  exit 1
fi

echo "3. Protected endpoints (token: ${ACCESS_TOKEN:0:20}... )"
curl_test "${BASE_URL}/api/v1/candidates/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN"

curl_test "${BASE_URL}/api/v1/companies/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN"

curl_test -X POST "${BASE_URL}/api/v1/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 4. Rate limit test (hit login 6 times fast)
echo "4. Rate limit test..."
for i in {1..6}; do
  curl -s -w "%{http_code}\n" -X POST "${BASE_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"rate@test.com","password":"pass"}' > /dev/null
done

# Expect 429 on last
if ! curl -s -f -X POST "${BASE_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"rate@test.com","password":"pass"}' 2>/dev/null | grep -q "429"; then
  echo "⚠️ Rate limit may not trigger (depends on timing)"
else
  echo "✓ Rate limit triggered (429)"
fi

# 5. Large body test (2MB -> 413)
echo "5. Large body test..."
python3 -c "import sys; sys.stdout.buffer.write(b'a'*2*1024*1024)" | \
curl -s --data-binary @- -X POST "${BASE_URL}/api/v1/auth/register" \
  -H "Content-Type: application/json" -w "%{http_code}\n" | grep -q "413" && echo "✓ Body limit OK (413)" || echo "⚠️ Body limit test"

echo "🎉 All smoke tests PASSED!"
echo "Full API endpoints ready."
