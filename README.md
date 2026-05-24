# seatbelt

seatbelt is a command line util to be used as VSCode Copilot hooks.

## Usage

### UserPromptSubmit hook

```bash
seatbelt scan-secrets --action=block  ## block prompt submission
```

This will scan the prompt before submission for secrets. If any secrets are found, the behavior will depend on the `action` specified.

- `--action=block`: blocks prompt submission with the following response to vscode: 
```json
{
  "continue": false,
  "stopReason": "Security policy violation",
  "systemMessage": "Secret from {location} in prompt detected"
}
```
`location` is the place where secrets originated, e.g. environment variables, `.env` file etc

#### how it detects secrets

seatbelt uses a combination of regex patterns and heuristics to identify potential secrets in the prompt. It looks for common patterns such as:
- API keys (e.g., `AIzaSy...`, `sk_live_...`)
- Database connection strings (e.g., `mongodb+srv://...`, `postgres://...`)
- Private keys (e.g., `-----BEGIN PRIVATE KEY-----` ...)
- Passwords (e.g., `password=...`, `pwd=...`)
- Other sensitive information
  - environment varables named 
    - `XXX_TOKEN` 
    - `XXX_KEY` 
  - well known config files that contains secrets:
    - `~/.m2/settings.xml`
    - `~/.pip/pip.conf`
    - `.env`
- Customized regex patterns defined in `~/.seatbelt.yaml`
