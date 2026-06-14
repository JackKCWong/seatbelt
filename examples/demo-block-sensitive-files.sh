#!/bin/bash
set -e

cd "$(dirname "$0")/.."

run_test() {
    local label="$1"
    local input="$2"
    echo ""
    echo "=== $label ==="
    echo "Input: $input"
    echo "$input" | ./seatbelt.exe block-sensitive-files
}

echo "=== Build ==="
go build -o seatbelt.exe .

run_test "Test: Normal file (should ALLOW)" \
    '{"tool_name":"read_file","tool_input":{"filePath":"src/normal.go"}}'

run_test "Test: .env file (should DENY)" \
    '{"tool_name":"read_file","tool_input":{"filePath":".env"}}'

run_test "Test: Nested .env file (should DENY)" \
    '{"tool_name":"read_file","tool_input":{"filePath":"config/prod/.env"}}'

run_test "Test: credentials file (should DENY)" \
    '{"tool_name":"read_file","tool_input":{"filePath":"credentials"}}'

run_test "Test: ~/.npmrc style path (should DENY)" \
    '{"tool_name":"read_file","tool_input":{"filePath":"/home/user/.npmrc"}}'

run_test "Test: Empty stdin (should ALLOW)" ''

echo ""
echo "=== Cleanup ==="
rm -f seatbelt.exe

echo ""
echo "Done!"
