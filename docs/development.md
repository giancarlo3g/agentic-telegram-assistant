# Development Guide

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **Git**: [Download Git](https://git-scm.com/)
- **Docker** (optional): [Download Docker](https://docker.com/)
- **Telegram Bot Token**: Get from [@BotFather](https://t.me/botfather)
- **OpenAI API Key**: Get from [OpenAI Platform](https://platform.openai.com/)
- **Google Calendar API Credentials**: Set up in [Google Cloud Console](https://console.cloud.google.com/)

### Development Environment Setup

#### 1. **Clone the Repository**
```bash
git clone <repository-url>
cd calendar-assistant-bot
```

#### 2. **Install Dependencies**
```bash
go mod download
```

#### 3. **Configure Environment**
```bash
cp env.example .env
# Edit .env with your credentials
```

#### 4. **Set Up Google Credentials**
```bash
mkdir -p credentials
# Place your google-credentials.json file in credentials/
chmod 600 credentials/google-credentials.json
```

#### 5. **Verify Setup**
```bash
go build ./cmd/bot
```

## 🏗️ Project Structure

### Directory Layout
```
calendar-assistant-bot/
├── cmd/
│   └── bot/                 # Main application entry point
│       └── main.go
├── pkg/                     # Core packages
│   ├── ai/                  # AI processing and orchestration
│   │   ├── agent.go         # Main AI agent
│   │   └── openai.go        # OpenAI API integration
│   ├── calendar/            # Google Calendar integration
│   │   └── calendar.go      # Calendar service
│   ├── config/              # Configuration management
│   │   └── config.go        # Config loading and validation
│   ├── database/            # Data persistence
│   │   └── database.go      # File-based database
│   ├── telegram/            # Telegram Bot API
│   │   └── bot.go           # Bot wrapper
│   └── types/               # Shared data structures
│       └── types.go         # Common types and interfaces
├── docs/                    # Documentation
├── credentials/             # Google API credentials (gitignored)
├── data/                    # Database files (gitignored)
├── logs/                    # Log files (gitignored)
├── Dockerfile               # Docker image definition
├── docker-compose.go.yml    # Docker Compose configuration
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── Makefile                 # Build and deployment commands
└── README.md                # Project overview
```

### Package Responsibilities

#### `cmd/bot/`
- **Purpose**: Application entry point and main orchestration
- **Responsibilities**: 
  - Initialize all services
  - Coordinate bot startup
  - Handle main event loop

#### `pkg/ai/`
- **Purpose**: AI processing and decision making
- **Responsibilities**:
  - Process user messages through OpenAI
  - Execute calendar actions
  - Manage conversation flow

#### `pkg/calendar/`
- **Purpose**: Google Calendar API integration
- **Responsibilities**:
  - CRUD operations on calendar events
  - Date parsing and validation
  - API error handling

#### `pkg/database/`
- **Purpose**: Data persistence and retrieval
- **Responsibilities**:
  - Store user interactions
  - Provide conversation context
  - Manage data files

#### `pkg/telegram/`
- **Purpose**: Telegram Bot API wrapper
- **Responsibilities**:
  - Send/receive messages
  - Handle bot updates
  - Message formatting

#### `pkg/config/`
- **Purpose**: Configuration management
- **Responsibilities**:
  - Load environment variables
  - Validate configuration
  - Provide secure access to secrets

#### `pkg/types/`
- **Purpose**: Shared data structures
- **Responsibilities**:
  - Define common interfaces
  - Provide type safety
  - Enable package communication

## 🔧 Development Workflow

### 1. **Feature Development**

#### Create a Feature Branch
```bash
git checkout -b feature/new-calendar-feature
```

#### Implement the Feature
- Follow Go coding standards
- Add comprehensive logging
- Include error handling
- Write tests for new functionality

#### Test Your Changes
```bash
# Run tests
go test ./...

# Build the application
go build ./cmd/bot

# Test locally
go run ./cmd/bot
```

#### Commit Your Changes
```bash
git add .
git commit -m "feat: add new calendar feature

- Implemented new calendar functionality
- Added comprehensive error handling
- Updated documentation"
```

### 2. **Code Review Process**

#### Before Submitting
- Ensure all tests pass
- Run `go fmt` and `go vet`
- Update documentation if needed
- Test with different configurations

#### Pull Request Template
```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Local testing completed
- [ ] All tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows Go standards
- [ ] Error handling implemented
- [ ] Logging added where appropriate
- [ ] Documentation updated
```

### 3. **Testing Strategy**

#### Unit Tests
```go
// pkg/ai/agent_test.go
func TestAgent_ProcessUserMessage(t *testing.T) {
    // Arrange
    mockOpenAI := &MockOpenAIService{}
    mockCalendar := &MockCalendarService{}
    mockTelegram := &MockTelegramBot{}
    mockDB := &MockDatabase{}
    
    agent := NewAgent(mockOpenAI, mockCalendar, mockTelegram, mockDB)
    
    // Act
    err := agent.ProcessUserMessage(123, 456, "Hello")
    
    // Assert
    assert.NoError(t, err)
    // Add more assertions
}
```

#### Integration Tests
```go
// tests/integration_test.go
func TestCalendarIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Test with real Google Calendar API
    // Use test calendar ID
}
```

#### Test Utilities
```go
// pkg/testing/mocks.go
type MockOpenAIService struct {
    ProcessMessageFunc func(userContext, userMessage string) (*types.AIResponse, error)
}

func (m *MockOpenAIService) ProcessMessage(userContext, userMessage string) (*types.AIResponse, error) {
    if m.ProcessMessageFunc != nil {
        return m.ProcessMessageFunc(userContext, userMessage)
    }
    return nil, errors.New("mock not implemented")
}
```

## 📝 Coding Standards

### Go Code Style

#### 1. **Formatting**
```bash
# Format code
go fmt ./...

# Check code style
go vet ./...

# Run linter
golangci-lint run
```

#### 2. **Naming Conventions**
```go
// Package names: lowercase, single word
package ai

// Function names: PascalCase for exported, camelCase for private
func ProcessUserMessage() {} // exported
func processMessage() {}     // private

// Variable names: camelCase
var userID int64
var calendarService *calendar.Service

// Constants: PascalCase
const MaxRetries = 3
const DefaultTimeout = 10 * time.Second
```

#### 3. **Error Handling**
```go
// Always check errors
if err != nil {
    return fmt.Errorf("failed to process message: %w", err)
}

// Use custom error types for specific errors
type CalendarError struct {
    Operation string
    Err       error
}

func (e *CalendarError) Error() string {
    return fmt.Sprintf("calendar operation %s failed: %v", e.Operation, e.Err)
}
```

#### 4. **Logging**
```go
// Use structured logging with context
log.Printf("Processing message from user %d: %s", userID, message)

// Include relevant IDs in all log entries
log.Printf("AI response for user %d: Action=%s, Message=%s", 
    userID, aiResponse.Action, aiResponse.Message)
```

### 5. **Documentation**
```go
// Package documentation
// Package ai provides AI processing capabilities for the calendar assistant bot.
package ai

// Function documentation
// ProcessUserMessage handles a complete user message flow including AI processing,
// action execution, and response generation.
func (a *Agent) ProcessUserMessage(userID, chatID int64, message string) error {
    // Implementation
}
```

## 🧪 Testing

### Running Tests

#### All Tests
```bash
go test ./...
```

#### Specific Package
```bash
go test ./pkg/ai
```

#### With Coverage
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Race Detection
```bash
go test -race ./...
```

#### Benchmark Tests
```bash
go test -bench=. ./...
```

### Test Structure

#### Test Files
```go
// pkg/ai/agent_test.go
package ai

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestAgent_NewAgent(t *testing.T) {
    // Test constructor
}

func TestAgent_ProcessUserMessage(t *testing.T) {
    // Test main functionality
}

func TestAgent_executeAIAction(t *testing.T) {
    // Test action execution
}
```

#### Mock Objects
```go
// pkg/ai/mocks_test.go
type MockOpenAIService struct {
    mock.Mock
}

func (m *MockOpenAIService) ProcessMessage(userContext, userMessage string) (*types.AIResponse, error) {
    args := m.Called(userContext, userMessage)
    return args.Get(0).(*types.AIResponse), args.Error(1)
}
```

## 🔍 Debugging

### Logging

#### Log Levels
```go
// Development: verbose logging
log.Printf("DEBUG: Processing message: %s", message)

// Production: essential logging only
log.Printf("INFO: Message processed for user %d", userID)
```

#### Context Information
```go
// Always include relevant context
log.Printf("Processing message from user %d (chatID %d): %s", 
    userID, chatID, message)
```

### Debug Mode

#### Environment Variable
```bash
DEBUG=true go run ./cmd/bot
```

#### Debug Functions
```go
func debugLog(format string, args ...interface{}) {
    if os.Getenv("DEBUG") == "true" {
        log.Printf("DEBUG: "+format, args...)
    }
}
```

### Error Tracing

#### Error Wrapping
```go
if err != nil {
    return fmt.Errorf("failed to create calendar service: %w", err)
}
```

#### Stack Traces
```go
import "runtime/debug"

func handlePanic() {
    if r := recover(); r != nil {
        log.Printf("Panic recovered: %v\n%s", r, debug.Stack())
    }
}
```

## 🚀 Performance Optimization

### 1. **Database Optimization**

#### Connection Pooling
```go
// For future database implementations
type DatabasePool struct {
    connections chan *Database
    maxConnections int
}

func (p *DatabasePool) GetConnection() *Database {
    select {
    case conn := <-p.connections:
        return conn
    default:
        return p.createNewConnection()
    }
}
```

#### Caching
```go
type Cache struct {
    data map[string]interface{}
    mu   sync.RWMutex
    ttl  time.Duration
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if item, exists := c.data[key]; exists {
        return item, true
    }
    return nil, false
}
```

### 2. **API Optimization**

#### Rate Limiting
```go
type RateLimiter struct {
    requests chan struct{}
    interval time.Duration
}

func (r *RateLimiter) Wait() {
    select {
    case r.requests <- struct{}{}:
        // Request allowed
    case <-time.After(r.interval):
        // Rate limit exceeded
    }
}
```

#### Request Batching
```go
type BatchProcessor struct {
    requests chan Request
    batchSize int
    timeout   time.Duration
}

func (b *BatchProcessor) Process() {
    var batch []Request
    timer := time.NewTimer(b.timeout)
    
    for {
        select {
        case req := <-b.requests:
            batch = append(batch, req)
            if len(batch) >= b.batchSize {
                b.processBatch(batch)
                batch = batch[:0]
            }
        case <-timer.C:
            if len(batch) > 0 {
                b.processBatch(batch)
                batch = batch[:0]
            }
            timer.Reset(b.timeout)
        }
    }
}
```

## 🔒 Security Best Practices

### 1. **Input Validation**

#### Sanitize User Input
```go
func sanitizeInput(input string) string {
    // Remove potentially dangerous characters
    input = strings.TrimSpace(input)
    input = html.EscapeString(input)
    return input
}
```

#### Validate Dates
```go
func validateDate(dateStr string) (time.Time, error) {
    // Check for SQL injection patterns
    if strings.Contains(dateStr, ";") || strings.Contains(dateStr, "--") {
        return time.Time{}, errors.New("invalid date format")
    }
    
    // Parse and validate date
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return time.Time{}, fmt.Errorf("invalid date format: %w", err)
    }
    
    return date, nil
}
```

### 2. **Authentication**

#### Token Validation
```go
func validateToken(token string) bool {
    // Check token format
    if !strings.HasPrefix(token, "sk-") {
        return false
    }
    
    // Check token length
    if len(token) < 20 {
        return false
    }
    
    return true
}
```

### 3. **Rate Limiting**

#### User Rate Limiting
```go
type UserRateLimiter struct {
    requests map[int64][]time.Time
    mu       sync.RWMutex
    limit    int
    window   time.Duration
}

func (r *UserRateLimiter) Allow(userID int64) bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    now := time.Now()
    userRequests := r.requests[userID]
    
    // Remove old requests outside the window
    var validRequests []time.Time
    for _, reqTime := range userRequests {
        if now.Sub(reqTime) < r.window {
            validRequests = append(validRequests, reqTime)
        }
    }
    
    if len(validRequests) >= r.limit {
        return false
    }
    
    r.requests[userID] = append(validRequests, now)
    return true
}
```

## 📚 Documentation

### Code Documentation

#### Function Documentation
```go
// ProcessUserMessage handles a complete user message flow including AI processing,
// action execution, and response generation. It retrieves user context from the
// database, sends the message to OpenAI for processing, stores the interaction,
// executes any calendar actions, and sends the response back to the user.
//
// Parameters:
//   - userID: Telegram user ID for context retrieval and storage
//   - chatID: Telegram chat ID for sending responses
//   - message: User's message text to process
//
// Returns:
//   - error: Any error that occurred during processing
//
// Example:
//   err := agent.ProcessUserMessage(123, 456, "What's on my calendar today?")
//   if err != nil {
//       log.Printf("Failed to process message: %v", err)
//   }
func (a *Agent) ProcessUserMessage(userID, chatID int64, message string) error {
    // Implementation
}
```

#### Package Documentation
```go
// Package ai provides AI processing capabilities for the calendar assistant bot.
// It includes the main AI agent that orchestrates user interactions, OpenAI
// integration for natural language processing, and action execution for calendar
// operations.
//
// The package is designed with dependency injection to enable easy testing and
// flexible service composition. All AI processing is stateless and thread-safe.
//
// Example usage:
//   openaiService := ai.NewOpenAIService(apiKey)
//   calendarService := calendar.NewService(calAPI, calendarID)
//   telegramBot := telegram.NewBot(token)
//   database := database.NewDatabase("./data")
//   
//   agent := ai.NewAgent(openaiService, calendarService, telegramBot, database)
//   err := agent.ProcessUserMessage(userID, chatID, message)
package ai
```

### API Documentation

#### OpenAPI/Swagger
```yaml
# docs/api/openapi.yaml
openapi: 3.0.0
info:
  title: Calendar Assistant Bot API
  version: 1.0.0
  description: API for the Calendar Assistant Bot

paths:
  /webhook/telegram:
    post:
      summary: Handle Telegram webhook
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/TelegramUpdate'
      responses:
        '200':
          description: Webhook processed successfully
```

## 🚀 Deployment

### Local Development

#### Hot Reload
```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Create .air.toml configuration
# Run with hot reload
air
```

#### Environment Switching
```bash
# Development
cp .env.development .env

# Staging
cp .env.staging .env

# Production
cp .env.production .env
```

### Docker Development

#### Development Container
```dockerfile
# Dockerfile.dev
FROM golang:1.21-alpine

WORKDIR /app

# Install development tools
RUN go install github.com/cosmtrek/air@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Run with hot reload
CMD ["air"]
```

#### Docker Compose Development
```yaml
# docker-compose.dev.yml
services:
  calendar-bot:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
      - ./credentials:/app/credentials:ro
    environment:
      - DEBUG=true
    ports:
      - "8080:8080"
```

### Production Deployment

#### Docker Production
```dockerfile
# Dockerfile.prod
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/bot

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/credentials ./credentials

CMD ["./main"]
```

#### Kubernetes Deployment
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: calendar-bot
spec:
  replicas: 3
  selector:
    matchLabels:
      app: calendar-bot
  template:
    metadata:
      labels:
        app: calendar-bot
    spec:
      containers:
      - name: calendar-bot
        image: calendar-assistant-bot:latest
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
        - name: GOOGLE_CREDENTIALS_FILE
          value: /app/credentials/google-credentials.json
        - name: GOOGLE_CALENDAR_ID
          valueFrom:
            secretKeyRef:
              name: bot-secrets
              key: google-calendar-id
        volumeMounts:
        - name: credentials
          mountPath: /app/credentials
          readOnly: true
      volumes:
      - name: credentials
        secret:
          secretName: google-credentials
```

## 🤝 Contributing

### Contribution Guidelines

#### 1. **Fork and Clone**
```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/yourusername/calendar-assistant-bot.git
cd calendar-assistant-bot

# Add upstream remote
git remote add upstream https://github.com/originalowner/calendar-assistant-bot.git
```

#### 2. **Create Feature Branch**
```bash
git checkout -b feature/your-feature-name
```

#### 3. **Make Changes**
- Follow the coding standards
- Add tests for new functionality
- Update documentation
- Ensure all tests pass

#### 4. **Commit and Push**
```bash
git add .
git commit -m "feat: add your feature description"
git push origin feature/your-feature-name
```

#### 5. **Create Pull Request**
- Use the provided PR template
- Describe your changes clearly
- Include any relevant issue numbers
- Request review from maintainers

### Code Review Process

#### Review Checklist
- [ ] Code follows Go standards
- [ ] Tests are included and pass
- [ ] Documentation is updated
- [ ] Error handling is appropriate
- [ ] Logging is comprehensive
- [ ] Security considerations are addressed

#### Review Comments
```go
// Good: Clear, actionable feedback
// Consider adding validation for empty event titles
if event.Title == "" {
    return errors.New("event title cannot be empty")
}

// Better: Suggest specific improvements
// Add validation to prevent empty event titles
if strings.TrimSpace(event.Title) == "" {
    return errors.New("event title cannot be empty or whitespace-only")
}
```

---

*Next: [Deployment](deployment.md) - Docker and production deployment*

