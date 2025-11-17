# GitReviewed 🔍

An automated code reviewer that scans GitHub Pull Requests for hardcoded secrets and sends alerts to Slack.

## Features

- Detects 18+ types of hardcoded secrets (AWS keys, GitHub tokens, API keys, etc.)
- Real-time Slack notifications with severity levels
- Automatic webhook verification for security
- Smart filtering to reduce false positives
- Beautiful formatted Slack messages with Block Kit

## Setup

### Prerequisites
- Go 1.21+
- GitHub account
- Slack workspace

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and fill in credentials
3. Build: `go build -o bin/reviewer cmd/reviewer/main.go`
4. Run: `./bin/reviewer`

### Configuration

Set these environment variables:
- `GITHUB_TOKEN`: GitHub Personal Access Token with `repo` scope
- `WEBHOOK_SECRET`: Random secret for webhook verification
- `SLACK_TOKEN`: Slack Bot Token (xoxb-...)
- `SLACK_CHANNEL`: Channel to post alerts (e.g., #code-reviews)

### GitHub Webhook Setup

1. Go to your repo → Settings → Webhooks
2. Add webhook with:
   - URL: `https://your-domain.com/webhook`
   - Content type: `application/json`
   - Secret: Your `WEBHOOK_SECRET`
   - Events: Pull requests only

## Detected Secret Types

- AWS Access Keys & Secret Keys
- GitHub Tokens (Personal, OAuth, App)
- OpenAI API Keys
- Slack Tokens & Webhooks
- Google API Keys
- Stripe API Keys
- Private Keys (RSA, DSA, EC)
- Database Connection Strings
- JWT Tokens
- And more...

## Endpoints

- `GET /health` - Health check
- `POST /webhook` - GitHub webhook endpoint
- `GET /test-slack` - Test Slack connection

## Development
```bash
# Run tests
go test ./...

# Build
make build

# Run locally
make run
```

## License

MIT