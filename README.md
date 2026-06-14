# seatbelt

seatbelt is a command line util to be used as VSCode Copilot hooks.

## Usage

### PreToolUse hook

#### block-sensitive-files during read_file

```bash
seatbelt block-sensitive-files
```

This will block VSCode agent from reading any sensitive files matching glob patterns in `.copilotdeny`.

Whenever the VSCode agent tries to read the files matching the patterns in `.copilotdeny`, the command outputs the following json to deny the access.

```json
{
  "continue": true,
  "hookSpecificOutput": {
    "hookEventName": "PreToolUse",
    "permissionDecision": "deny",
    "permissionDecisionReason": "sensitive file read blocked by policy"
  }
}
```
