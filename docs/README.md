# Calendar Assistant Bot - Technical Documentation

Welcome to the technical documentation for the Calendar Assistant Bot, a Go-based Telegram bot that integrates with OpenAI and Google Calendar APIs.

## 📚 Documentation Structure

- **[Getting Started](getting-started.md)** - Quick setup and first run
- **[Architecture Overview](architecture.md)** - High-level system design
- **[API Reference](api-reference.md)** - Complete function and type documentation
- **[Configuration](configuration.md)** - Environment variables and settings
- **[Development Guide](development.md)** - Contributing and extending the bot
- **[Deployment](deployment.md)** - Docker and production deployment

## 🚀 Quick Start

```bash
# Clone and setup
git clone <repository>
cd calendar-assistant-bot
cp env.example .env
# Edit .env with your credentials

# Run locally
go run ./cmd/bot

# Or with Docker
docker-compose -f docker-compose.go.yml up -d
```

## 🔧 Prerequisites

- Go 1.21+
- Telegram Bot Token
- OpenAI API Key
- Google Calendar API Credentials
- Docker (optional)

## 📖 What You'll Learn

This documentation covers:
- Complete codebase architecture and design patterns
- All Go types, interfaces, and functions with examples
- API integration patterns and error handling
- Database design and persistence strategies
- Testing and debugging approaches
- Production deployment considerations

---

*For questions or contributions, please refer to the [Development Guide](development.md).*

