# Calendar Assistant Bot

A Telegram bot that helps manage Google Calendar events using AI-powered natural language processing.

## 🏗️ Project Structure

The project is organized into a clean, modular Go package structure:

```
go-telegram-ai/
├── cmd/
│   └── bot/                 # Main application entry point
│       └── main.go
├── pkg/                     # Reusable packages
│   ├── ai/                  # AI-related functionality
│   │   ├── agent.go         # AI agent coordination
│   │   └── openai.go        # OpenAI API integration
│   ├── calendar/            # Google Calendar operations
│   │   └── calendar.go      # Calendar service
│   ├── config/              # Configuration management
│   │   └── config.go        # App configuration
│   ├── database/            # Data persistence
│   │   └── database.go      # Interaction storage
│   ├── telegram/            # Telegram bot functionality
│   │   └── bot.go           # Bot operations
│   └── types/               # Shared data structures
│       └── types.go         # Common types
├── credentials/              # Google service account credentials
├── data/                     # Database storage directory
├── logs/                     # Application logs
├── Dockerfile                # Docker container definition
├── docker-compose.go.yml     # Docker Compose configuration
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── Makefile                  # Build and deployment commands
```

## 🚀 Features

- **Natural Language Processing**: Understand calendar requests in plain English
- **Google Calendar Integration**: Create, read, update, and delete events
- **AI-Powered Responses**: Uses OpenAI GPT-4o-mini for intelligent responses
- **Conversation Memory**: Remembers user context and interaction history
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **Modular Architecture**: Clean separation of concerns for maintainability

## 🛠️ Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Google Cloud Platform account with Calendar API enabled
- OpenAI API key
- Telegram Bot Token

## 📋 Setup

### 1. Environment Variables

Create a `.env` file or set environment variables:

```bash
TELEGRAM_TOKEN=your_telegram_bot_token
OPENAI_API_KEY=your_openai_api_key
GOOGLE_CREDENTIALS_FILE=/app/credentials/google-credentials.json
GOOGLE_CALENDAR_ID=your_calendar_id@group.calendar.google.com
PORT=8080
```

### 2. Google Calendar Setup

1. Enable Google Calendar API in Google Cloud Console
2. Create a service account and download credentials JSON
3. Place the credentials file in `credentials/google-credentials.json`
4. Share your calendar with the service account email

### 3. Build and Run

#### Using Docker (Recommended)

```bash
# Build the Docker image
make docker-build

# Run the bot
make docker-run

# View logs
make view-logs

# Stop the bot
make docker-stop
```

#### Using Go directly

```bash
# Build the binary
go build -o calendar-bot ./cmd/bot

# Run the bot
./calendar-bot
```

## 💬 Usage Examples

The bot understands natural language requests:

- "Get events for today"
- "What's on my calendar tomorrow?"
- "Create a meeting with John tomorrow at 2pm"
- "Delete the meeting at 3pm today"
- "What time is it?"

## 🔧 Development

### Adding New Features

1. **New Calendar Operations**: Add methods to `pkg/calendar/calendar.go`
2. **New AI Actions**: Update the system prompt in `pkg/ai/openai.go`
3. **New Bot Commands**: Extend `pkg/telegram/bot.go`
4. **Data Models**: Add types to `pkg/types/types.go`

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Code Organization Principles

- **Single Responsibility**: Each package has one clear purpose
- **Dependency Injection**: Services are injected rather than created internally
- **Interface Segregation**: Small, focused interfaces
- **Clean Architecture**: Business logic separated from infrastructure

## 📦 Package Details

### `pkg/ai`
- **agent.go**: Coordinates between all services and executes AI decisions
- **openai.go**: Handles OpenAI API communication and response parsing

### `pkg/calendar`
- **calendar.go**: Google Calendar API operations (CRUD events)

### `pkg/config`
- **config.go**: Environment variable loading and validation

### `pkg/database`
- **database.go**: JSON-based interaction storage with thread safety

### `pkg/telegram`
- **bot.go**: Telegram Bot API integration and message handling

### `pkg/types`
- **types.go**: Shared data structures used across packages

## 🚀 Deployment

### Docker Deployment

```bash
# Build and run
make docker-build
make docker-run

# Production deployment
docker-compose -f docker-compose.go.yml up -d
```

### Kubernetes Deployment

The modular structure makes it easy to deploy individual components as separate services in Kubernetes.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following the package structure
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Database Deadlocks**: Fixed in the current version with proper mutex handling
2. **Google Calendar Timeouts**: Added 10-second timeouts to all API calls
3. **Import Path Issues**: Ensure all packages use the correct module path

### Debug Mode

Enable detailed logging by setting the log level in the configuration.

## 🔮 Future Enhancements

- [ ] Webhook support for Telegram
- [ ] Multiple calendar support
- [ ] Event reminders and notifications
- [ ] Calendar sharing and collaboration
- [ ] REST API for external integrations
- [ ] Metrics and monitoring
- [ ] Multi-language support
