#!/usr/bin/env bash

BASE_URL="http://localhost:8080"
JWT_SECRET="secret123"

USER_UUID="11111111-1111-1111-1111-111111111111"
BOOKABLE_UUID="a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

b64url() {
  openssl base64 -e -A | tr '+/' '-_' | tr -d '='
}

HEADER=$(echo -n '{"alg":"HS256","typ":"JWT"}' | b64url)
PAYLOAD=$(echo -n "{\"user_id\":\"${USER_UUID}\",\"role\":\"ADMIN\",\"exp\":1893456000}" | b64url)
SIGNATURE=$(echo -n "${HEADER}.${PAYLOAD}" | openssl dgst -binary -sha256 -hmac "$JWT_SECRET" | b64url)
TOKEN="${HEADER}.${PAYLOAD}.${SIGNATURE}"

echo "=== Testing POST /api/v1/bookings ==="
curl -i -X POST "${BASE_URL}/api/v1/bookings" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"bookable_id\": \"${BOOKABLE_UUID}\",
    \"user_id\": \"${USER_UUID}\",
    \"start_time\": \"2026-09-04T10:00:00Z\",
    \"end_time\": \"2026-09-04T11:00:00Z\"
  }"