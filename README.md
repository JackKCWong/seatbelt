# seatbelt

seatbelt is a command line util to be used as VSCode Copilot hooks.

## Usage

### installation

```bash
go install github.com/JackKCWong/seatbelt@latest
```

### PreToolUse hook

Add the following config to `.github/hooks/`, e.g. [.github/hooks/seatbelt.json](.github/hooks/seatbelt.json)

```json
{
	"hooks": {
		"PreToolUse": [
			{
				"type": "command",
				"command": "seatbelt block-sensitive-files"
			}
		]
	}
}
```

and a `.copilotdeny` file in your current respository root.

This will block VSCode agent from reading any sensitive files matching glob patterns in `.copilotdeny`.

***How it works:***


It reads the hook input from stdin and output the result to stdout to tell Copilot whether the file can be read. e.g. 

```bash
echo '{"tool_name":"read_file","tool_input":{"filePath":".env"}}' | seatbelt block-sensitive-files
```

```json
{
  "continue": true,
  "hookSpecificOutput": {
    "additionalContext": "Access to .env was blocked by sensitive file policy (matched pattern: **/.env).",
    "hookEventName": "PreToolUse",
    "permissionDecision": "deny",
    "permissionDecisionReason": "sensitive file read blocked by policy"
  }
}
```

***Note: it CANNOT block you from attaching the file to via the chat window because that does not invoke a tool use.***
