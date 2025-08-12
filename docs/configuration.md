# Configuration

## 🔧 Environment Variables

The Calendar Assistant Bot uses environment variables for all configuration. This approach provides flexibility across different deployment environments while keeping sensitive data secure.

### Required Environment Variables

#### `TELEGRAM_TOKEN`
**Description**: Your Telegram bot token from [@BotFather](https://t.me/botfather)

**Format**: `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`

**Example**:
```bash
TELEGRAM_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
```

**How to get it**:
1. Message [@BotFather](https://t.me/botfather) on Telegram
2. Send `/newbot` command
3. Follow the prompts to create your bot
4. Copy the token provided

#### `OPENAI_API_KEY`
**Description**: Your OpenAI API key for GPT-4o-mini access

**Format**: `sk-...` (starts with `sk-`)

**Example**:
```bash
OPENAI_API_KEY=sk-1234567890abcdefghijklmnopqrstuvwxyz
```

**How to get it**:
1. Visit [OpenAI Platform](https://platform.openai.com/api-keys)
2. Sign in or create an account
3. Click "Create new secret key"
4. Copy the generated key

#### `GOOGLE_CREDENTIALS_FILE`
**Description**: Path to your Google service account credentials JSON file

**Format**: File path relative to the application root

**Example**:
```bash
GOOGLE_CREDENTIALS_FILE=./credentials/google-credentials.json
```

**How to get it**:
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google Calendar API
4. Create a service account
5. Download the JSON credentials file
6. Place it in your `credentials/` directory

#### `GOOGLE_CALENDAR_ID`
**Description**: ID of the Google Calendar to manage

**Format**: Usually `primary` for main calendar, or specific calendar ID

**Example**:
```bash
GOOGLE_CALENDAR_ID=primary
# or
GOOGLE_CALENDAR_ID=abc123@group.calendar.google.com
```

**How to get it**:
1. Go to [Google Calendar](https://calendar.google.com/)
2. Find your calendar in the left sidebar
3. Click the three dots next to the calendar name
4. Select "Settings and sharing"
5. Copy the "Calendar ID" (usually `primary` for main calendar)

### Optional Environment Variables

#### `PORT`
**Description**: HTTP server port (if implementing webhooks)

**Default**: `8080`

**Example**:
```bash
PORT=3000
```

## 📁 Configuration Files

### `.env` File

Create a `.env` file in your project root with all required variables:

```bash
# Telegram Configuration
TELEGRAM_TOKEN=your_telegram_bot_token_here

# OpenAI Configuration
OPENAI_API_KEY=your_openai_api_key_here

# Google Calendar Configuration
GOOGLE_CREDENTIALS_FILE=./credentials/google-credentials.json
GOOGLE_CALENDAR_ID=primary

# Optional Configuration
PORT=8080
```

### `env.example`

The project includes an `env.example` file as a template:

```bash
# Copy this file to .env and fill in your values
TELEGRAM_TOKEN=
OPENAI_API_KEY=
GOOGLE_CREDENTIALS_FILE=./credentials/google-credentials.json
GOOGLE_CALENDAR_ID=primary
PORT=8080
```

## 🔐 Security Considerations

### Credential Management

#### 1. **Never Commit Credentials**
```bash
# .gitignore should contain:
.env
credentials/
```

#### 2. **File Permissions**
```bash
# Restrict access to credentials directory
chmod 600 credentials/google-credentials.json
chmod 700 credentials/
```

#### 3. **Environment Variable Security**
- Use `.env` files for local development
- Use secure environment variable injection in production
- Rotate API keys regularly
- Monitor API usage for anomalies

### Production Deployment

#### Docker Compose
```yaml
services:
  calendar-bot:
    environment:
      - TELEGRAM_TOKEN=${TELEGRAM_TOKEN}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - GOOGLE_CREDENTIALS_FILE=/app/credentials/google-credentials.json
      - GOOGLE_CALENDAR_ID=${GOOGLE_CALENDAR_ID}
```

#### Kubernetes
```yaml
env:
- name: TELEGRAM_TOKEN
  valueFrom:
    secretKeyRef:
      name: bot-secrets
      key: telegram-token
- name: OPENAI_API_KEY
  valueFrom:
    secretKeyRef:
      name: bot-secrets
      key: openai-api-key
```

## 🏗️ Configuration Loading

### Go Code Implementation

The configuration is loaded in `pkg/config/config.go`:

```go
func Load() (*Config, error) {
    // Load .env file if it exists
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found: %v", err)
    }
    
    config := &Config{
        TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
        OpenAIKey:     os.Getenv("OPENAI_API_KEY"),
        GoogleCreds:   os.Getenv("GOOGLE_CREDENTIALS_FILE"),
        CalendarID:    os.Getenv("GOOGLE_CALENDAR_ID"),
        Port:          getEnvWithDefault("PORT", "8080"),
    }
    
    // Validate required fields
    if err := validateConfig(config); err != nil {
        return nil, err
    }
    
    return config, nil
}
```

### Validation

The configuration loader validates all required fields:

```go
func validateConfig(config *Config) error {
    if config.TelegramToken == "" {
        return &ConfigError{Field: "TELEGRAM_TOKEN", Err: errors.New("required")}
    }
    if config.OpenAIKey == "" {
        return &ConfigError{Field: "OPENAI_API_KEY", Err: errors.New("required")}
    }
    if config.GoogleCreds == "" {
        return &ConfigError{Field: "GOOGLE_CREDENTIALS_FILE", Err: errors.New("required")}
    }
    if config.CalendarID == "" {
        return &ConfigError{Field: "GOOGLE_CALENDAR_ID", Err: errors.New("required")}
    }
    return nil
}
```

## 🔍 Configuration Debugging

### Logging Configuration

The bot logs configuration details (with sensitive data masked) on startup:

```go
log.Printf("Creating bot with config: Telegram=%s, OpenAI=%s, GoogleCreds=%s, CalendarID=%s",
    config.MaskToken(cfg.TelegramToken),
    config.MaskToken(cfg.OpenAIKey),
    cfg.GoogleCreds,
    cfg.CalendarID)
```

### Token Masking

Sensitive tokens are masked in logs for security:

```go
func MaskToken(token string) string {
    if len(token) <= 8 {
        return "***"
    }
    return token[:4] + "..." + token[len(token)-4:]
}
```

**Example output**:
```
2024/01/15 10:30:45 Creating bot with config: Telegram=1234...abcd, OpenAI=sk-1...xyz, GoogleCreds=./credentials/google-credentials.json, CalendarID=primary
```

## 🚀 Environment-Specific Configurations

### Development Environment

```bash
# .env.development
TELEGRAM_TOKEN=dev_bot_token
OPENAI_API_KEY=sk-dev-key
GOOGLE_CREDENTIALS_FILE=./credentials/dev-credentials.json
GOOGLE_CALENDAR_ID=dev-calendar-id
PORT=8080
```

### Staging Environment

```bash
# .env.staging
TELEGRAM_TOKEN=staging_bot_token
OPENAI_API_KEY=sk-staging-key
GOOGLE_CREDENTIALS_FILE=./credentials/staging-credentials.json
GOOGLE_CALENDAR_ID=staging-calendar-id
PORT=8080
```

### Production Environment

```bash
# .env.production
TELEGRAM_TOKEN=prod_bot_token
OPENAI_API_KEY=sk-prod-key
GOOGLE_CREDENTIALS_FILE=./credentials/prod-credentials.json
GOOGLE_CALENDAR_ID=prod-calendar-id
PORT=80
```

## 🔧 Configuration Testing

### Validation Script

Create a simple validation script to test your configuration:

```bash
#!/bin/bash
# validate-config.sh

echo "Validating Calendar Assistant Bot configuration..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ .env file not found"
    exit 1
fi

# Load environment variables
source .env

# Check required variables
required_vars=("TELEGRAM_TOKEN" "OPENAI_API_KEY" "GOOGLE_CREDENTIALS_FILE" "GOOGLE_CALENDAR_ID")

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "❌ $var is not set"
        exit 1
    else
        echo "✅ $var is set"
    fi
done

# Check if credentials file exists
if [ ! -f "$GOOGLE_CREDENTIALS_FILE" ]; then
    echo "❌ Google credentials file not found: $GOOGLE_CREDENTIALS_FILE"
    exit 1
else
    echo "✅ Google credentials file exists"
fi

echo "🎉 Configuration validation passed!"
```

### Usage

```bash
chmod +x validate-config.sh
./validate-config.sh
```

## 📊 Configuration Monitoring

### Health Checks

Monitor configuration health in production:

```go
func (c *Config) HealthCheck() error {
    // Test OpenAI API
    if err := testOpenAIAPI(c.OpenAIKey); err != nil {
        return fmt.Errorf("OpenAI API health check failed: %w", err)
    }
    
    // Test Google Calendar API
    if err := testGoogleCalendarAPI(c.GoogleCreds, c.CalendarID); err != nil {
        return fmt.Errorf("Google Calendar API health check failed: %w", err)
    }
    
    // Test Telegram Bot API
    if err := testTelegramAPI(c.TelegramToken); err != nil {
        return fmt.Errorf("Telegram API health check failed: %w", err)
    }
    
    return nil
}
```

### Metrics

Track configuration-related metrics:

```go
type ConfigMetrics struct {
    LastValidationTime time.Time
    ValidationErrors   int
    ConfigReloads      int
}
```

---

*Next: [Development Guide](development.md) - Contributing and extending the bot*

