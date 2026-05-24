
# VSCode Copilot hooks reference

## Hook input and output
Hooks communicate with VS Code through stdin (input) and stdout (output) using JSON.

Common input fields
Every hook receives a JSON object via stdin with these common fields:

```json
{
  "timestamp": "2026-02-09T10:30:00.000Z",
  "cwd": "/path/to/workspace",
  "sessionId": "session-identifier",
  "hookEventName": "PreToolUse",
  "transcript_path": "/path/to/transcript.json"
}
```

Common output format
Hooks can return JSON via stdout to influence agent behavior. All hooks support these output fields:

```json
{
  "continue": true,
  "stopReason": "Security policy violation",
  "systemMessage": "Unit tests failed"
}
```

## UserPromptSubmit
The UserPromptSubmit hook fires when the user submits a prompt.

UserPromptSubmit input
In addition to the common fields, UserPromptSubmit hooks receive a `prompt` field with the text the user submitted.

The UserPromptSubmit hook uses the common output format only.
