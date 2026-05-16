# sendr-cli

A command-line tool for the [Sendr](https://github.com/muhammedfazall/Sendr) transactional email API. Send emails, manage API keys, and configure your account directly from the terminal.

## Installation

**Requirements:** Go 1.25+

```bash
go install github.com/muhammedfazall/sendr-cli@latest
```

Or build from source:

```bash
git clone https://github.com/muhammedfazall/sendr-cli
cd sendr-cli
go build -o sendr main.go
```

---

## Quick Start

```bash
# 1. Point CLI to your Sendr backend (skip if using the hosted default)
sendr config set-url https://your-backend.com

# 2. Authenticate
sendr login

# 3. Create an API key
sendr keys create production

# 4. Send an email
sendr send --to user@example.com --subject "Hello" --body "Hello from sendr-cli"
```

---

## Commands

### Version

```bash
sendr --version
```

Prints the CLI version embedded during release builds.

---

### Authentication

#### `sendr login`
Opens your browser to authenticate with your Sendr account via Google OAuth. Stores credentials locally on success.

```bash
sendr login
```

#### `sendr logout`
Clears stored credentials from your local config.

```bash
sendr logout
```

---

### Sending Email

#### `sendr send`
Queues an email for delivery through the Sendr pipeline.

```bash
sendr send --to <email> --subject <subject> --body <body>
```

**Flags:**

| Flag | Required | Description |
|------|----------|-------------|
| `--to` | Yes | Recipient email address |
| `--subject` | Yes | Email subject line |
| `--body` | Yes | Email body (plain text), or `@path/to/file.txt` |
| `--no-wait` | No | Return after queuing instead of polling delivery status |

**Example:**

```bash
sendr send \
  --to user@example.com \
  --subject "Welcome to our service" \
  --body "Thanks for signing up!"
```

---

### API Keys

#### `sendr keys create <name>`
Creates a new API key and saves it to your local config. The full key is shown once — it is stored automatically.

```bash
sendr keys create production
```

#### `sendr keys list`
Lists all active API keys on your account.

```bash
sendr keys list
```

#### `sendr keys revoke <id>`
Revokes an API key by its ID. Use `sendr keys list` to find the ID.

```bash
sendr keys revoke <key-id>
```

---

### Configuration

#### `sendr config set-url <url>`
Sets the Sendr backend API URL. Defaults to `https://api.sendr.app`. Use this if you are self-hosting Sendr.

```bash
sendr config set-url https://your-backend.com
```

#### `sendr config view`
Shows your current CLI configuration — API URL, login status, and active API key prefix.

```bash
sendr config view
```

---

## How It Works

- `sendr login` starts a local HTTP server, opens your browser to the backend OAuth URL, and captures the token when Google redirects back.
- Credentials are stored in `~/.config/sendr/config.json`.
- `sendr logout` clears the login token but keeps the saved API key for future sends.
- `sendr send` authenticates requests using your stored API key.
- Emails are queued in the backend and delivered asynchronously via SendGrid with automatic retries.

---

## Config File

Stored at `~/.config/sendr/config.json`:

```json
{
  "api_url": "https://api.sendr.app",
  "token": "<jwt>",
  "api_key": "mk_live_..."
}
```

- `token` — JWT used for account-level operations (key management)
- `api_key` — used for sending emails, persists across logout/login

---

## Requirements

- A running [Sendr](https://github.com/muhammedfazall/Sendr) backend instance
- A Google account for authentication

---

## License

MIT
