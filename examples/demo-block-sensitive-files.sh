#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "=== Build ==="
go build -o seatbelt.exe .

echo ""
echo "=== Test: Normal file (should ALLOW) ==="
echo '{"tool_name":"read_file","tool_input":{"filePath":"src/normal.go","startLine":1,"endLine":200}}' | ./seatbelt.exe block-sensitive-files

echo ""
echo "=== Test: .env file (should DENY) ==="
echo '{"tool_name":"read_file","tool_input":{"filePath":".env","startLine":1,"endLine":200}}' | ./seatbelt.exe block-sensitive-files

echo ""
echo "=== Test: Nested .env file (should DENY) ==="
echo '{"tool_name":"read_file","tool_input":{"filePath":"config/prod/.env","startLine":1,"endLine":200}}' | ./seatbelt.exe block-sensitive-files

echo ""
echo "=== Test: credentials file (should DENY) ==="
echo '{"tool_name":"read_file","tool_input":{"filePath":"credentials","startLine":1,"endLine":200}}' | ./seatbelt.exe block-sensitive-files

echo ""
echo "=== Test: ~/.npmrc style path (should DENY) ==="
echo '{"tool_name":"read_file","tool_input":{"filePath":"/home/user/.npmrc","startLine":1,"endLine":200}}' | ./seatbelt.exe block-sensitive-files

echo ""
echo "=== Test: Empty stdin (should ALLOW) ==="
echo "" | ./seatbelt.exe block-sensitive-files

echo ""
echo "=== Cleanup ==="
rm -f seatbelt.exe

echo ""
echo "Done!"
