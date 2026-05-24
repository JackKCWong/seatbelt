#!/bin/bash
set -e

SEATBELT="./seatbelt.exe"

CLEAN_INPUT='{"timestamp":"2026-02-09T10:30:00.000Z","cwd":"/test","sessionId":"test-session","hookEventName":"UserPromptSubmit","prompt":"hello world"}'
echo "Testing seatbelt.exe with clean input..."
echo "Input: $CLEAN_INPUT"
OUTPUT=$(printf "$CLEAN_INPUT" | "$SEATBELT" scan-secrets)
EXIT_CODE=$?
echo "Output: $OUTPUT"
echo "Exit code: $EXIT_CODE"

if echo "$OUTPUT" | grep -q '"continue":true'; then
    echo "PASS: Clean input passed through"
else
    echo "FAIL: Expected continue:true"
    exit 1
fi

echo ""
SECRET_INPUT='{"timestamp":"2026-02-09T10:30:00.000Z","cwd":"/test","sessionId":"test-session","hookEventName":"UserPromptSubmit","prompt":"password=supersecret"}'
echo "Testing seatbelt.exe with secret in prompt..."
echo "Input: $SECRET_INPUT"
SECRET_OUTPUT=$(printf "$SECRET_INPUT" | "$SEATBELT" scan-secrets)
SECRET_EXIT=$?
echo "Output: $SECRET_OUTPUT"
echo "Exit code: $SECRET_EXIT"

if echo "$SECRET_OUTPUT" | grep -q '"continue":false'; then
    echo "PASS: Secret detected"
else
    echo "FAIL: Expected secret to be detected"
    exit 1
fi

echo ""

echo ""
SECRET_INPUT='{"timestamp":"2026-02-09T10:30:00.000Z","cwd":"/test","sessionId":"test-session","hookEventName":"UserPromptSubmit","prompt":"define API_KEY=supersecret to block this prompt"}'
echo "Testing seatbelt.exe with secret from environment..."
echo "Input: $SECRET_INPUT"
SECRET_OUTPUT=$(printf "$SECRET_INPUT" | "$SEATBELT" scan-secrets)
SECRET_EXIT=$?
echo "Output: $SECRET_OUTPUT"
echo "Exit code: $SECRET_EXIT"

if echo "$SECRET_OUTPUT" | grep -q '"continue":false'; then
    echo "PASS: Secret detected"
else
    echo "FAIL: Expected secret to be detected"
    exit 1
fi

echo ""
echo "All smoke tests passed!"