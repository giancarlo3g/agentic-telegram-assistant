# API Reference

## 📋 Table of Contents

- [Types](#types)
- [AI Package](#ai-package)
- [Calendar Package](#calendar-package)
- [Database Package](#database-package)
- [Telegram Package](#telegram-package)
- [Config Package](#config-package)
- [Main Application](#main-application)

---

## 🏷️ Types

### Core Data Structures

#### `AIResponse`
Represents the AI agent's response to a user message.

```go
type AIResponse struct {
    Action     string      `json:"action"`           // Action to perform
    Message    string      `json:"message"`          // Response message
    EventID    string      `json:"event_id,omitempty"`      // Event ID for updates/deletes
    EventTitle string      `json:"event_title,omitempty"`   // Event title for creation/updates
    EventDate  string      `json:"event_date,omitempty"`    // Event date (YYYY-MM-DD or relative)
    EventTime  string      `json:"event_time,omitempty"`    // Event time (HH:MM)
    EventDesc  string      `json:"event_description,omitempty"` // Event description
    EventLoc   string      `json:"event_location,omitempty"`    // Event location
    Actions    []AIAction  `json:"actions,omitempty"`      // Multiple actions for complex requests
}
```

**Fields:**
- `Action`: One of `"getEvents"`, `"makeEvent"`, `"updtEvent"`, `"delEvents"`, `"message"`, `"None"`
- `Message`: Human-readable response to the user
- `EventID`: Required for update/delete operations
- `EventTitle/EventDate/EventTime/EventDesc/EventLoc`: Required for create/update operations
- `Actions`: Array of actions for complex multi-step requests

#### `AIAction`
Represents a single action the AI wants to perform.

```go
type AIAction struct {
    Action     string `json:"action"`           // Action type
    EventDate  string `json:"event_date,omitempty"`    // Date for the action
    EventTitle string `json:"event_title,omitempty"`   // Event title
    EventTime  string `json:"event_time,omitempty"`    // Event time
    EventDesc  string `json:"event_description,omitempty"` // Event description
    EventLoc   string `json:"event_location,omitempty"`    // Event location
}
```

#### `CalendarEvent`
Represents a calendar event with all its properties.

```go
type CalendarEvent struct {
    ID          string    `json:"id"`           // Unique event identifier
    Summary     string    `json:"summary"`      // Event title/summary
    Description string    `json:"description"`  // Event description
    Start       time.Time `json:"start"`        // Event start time
    End         time.Time `json:"end"`          // Event end time
    Location    string    `json:"location"`     // Event location
}
```

#### `Interaction`
Represents a single user-AI interaction for context tracking.

```go
type Interaction struct {
    UserID      int64     `json:"user_id"`      // Telegram user ID
    Timestamp   time.Time `json:"timestamp"`    // When the interaction occurred
    UserMessage string    `json:"user_message"` // What the user said
    AIResponse  string    `json:"ai_response"`  // What the AI responded
    Action      string    `json:"action,omitempty"` // What action was taken
}
```

---

## 🤖 AI Package

### `pkg/ai/agent.go`

#### `Agent`
Main AI agent that orchestrates all interactions.

```go
type Agent struct {
    openaiService   *OpenAIService
    calendarService *calendar.Service
    telegramBot     *telegram.Bot
    database        *database.Database
}
```

#### `NewAgent()`
Creates a new AI agent instance.

```go
func NewAgent(
    openaiService *OpenAIService,
    calendarService *calendar.Service,
    telegramBot *telegram.Bot,
    database *database.Database,
) *Agent
```

**Parameters:**
- `openaiService`: OpenAI API service
- `calendarService`: Google Calendar service
- `telegramBot`: Telegram bot instance
- `database`: Database for storing interactions

**Returns:** `*Agent` - New agent instance

#### `ProcessUserMessage()`
Main entry point for processing user messages.

```go
func (a *Agent) ProcessUserMessage(userID, chatID int64, message string) error
```

**Parameters:**
- `userID`: Telegram user ID
- `chatID`: Telegram chat ID
- `message`: User's message text

**Returns:** `error` - Any error that occurred during processing

**Flow:**
1. Retrieves user context from database
2. Sends message to OpenAI for processing
3. Stores interaction in database
4. Executes AI's decision (calendar actions, etc.)
5. Sends response to user

#### `executeAIAction()`
Executes the action determined by the AI.

```go
func (a *Agent) executeAIAction(userID int64, aiResp types.AIResponse) (string, error)
```

**Parameters:**
- `userID`: Telegram user ID
- `aiResp`: AI's response with action details

**Returns:** `(string, error)` - Response message and any error

**Supported Actions:**
- `getEvents`: Retrieves events for a specific date
- `makeEvent`: Creates a new calendar event
- `updtEvent`: Updates an existing event
- `delEvents`: Deletes an event
- `message`: Sends a conversational response
- `None`: No action needed

### `pkg/ai/openai.go`

#### `OpenAIService`
Handles all OpenAI API interactions.

```go
type OpenAIService struct {
    client *openai.Client
}
```

#### `NewOpenAIService()`
Creates a new OpenAI service instance.

```go
func NewOpenAIService(apiKey string) *OpenAIService
```

**Parameters:**
- `apiKey`: OpenAI API key

**Returns:** `*OpenAIService` - New service instance

#### `ProcessMessage()`
Sends a message to OpenAI and parses the response.

```go
func (s *OpenAIService) ProcessMessage(userContext, userMessage string) (*types.AIResponse, error)
```

**Parameters:**
- `userContext`: Previous conversation context
- `userMessage`: Current user message

**Returns:** `(*types.AIResponse, error)` - Parsed AI response and any error

**System Prompt:**
The AI is instructed to:
- Handle calendar management tasks
- Provide conversational responses
- Use relative dates (today, tomorrow, yesterday)
- Return structured JSON responses
- Handle complex multi-step requests

---

## 📅 Calendar Package

### `pkg/calendar/calendar.go`

#### `Service`
Google Calendar service wrapper.

```go
type Service struct {
    service    *calapi.Service
    calendarID string
}
```

#### `NewService()`
Creates a new calendar service instance.

```go
func NewService(service *calapi.Service, calendarID string) *Service
```

**Parameters:**
- `service`: Google Calendar API service
- `calendarID`: Target calendar ID

**Returns:** `*Service` - New calendar service instance

#### `GetEvents()`
Retrieves events for a specific date.

```go
func (s *Service) GetEvents(dateStr string) ([]types.CalendarEvent, error)
```

**Parameters:**
- `dateStr`: Date string (YYYY-MM-DD, "today", "tomorrow", "yesterday")

**Returns:** `([]types.CalendarEvent, error)` - Events and any error

**Date Handling:**
- `"today"`: Current date
- `"tomorrow"`: Next day
- `"yesterday"`: Previous day
- `"YYYY-MM-DD"`: Specific date

**API Calls:**
- Uses `context.WithTimeout(10s)` for robustness
- Formats dates using `time.RFC3339`
- Orders events by start time

#### `CreateEvent()`
Creates a new calendar event.

```go
func (s *Service) CreateEvent(event types.CalendarEvent) error
```

**Parameters:**
- `event`: Event details to create

**Returns:** `error` - Any error that occurred

**Event Creation:**
- Converts internal types to Google Calendar API types
- Sets default duration to 1 hour if not specified
- Handles timezone conversion

#### `UpdateEvent()`
Updates an existing calendar event.

```go
func (s *Service) UpdateEvent(eventID string, event types.CalendarEvent) error
```

**Parameters:**
- `eventID`: ID of event to update
- `event`: New event details

**Returns:** `error` - Any error that occurred

#### `DeleteEvent()`
Deletes a calendar event.

```go
func (s *Service) DeleteEvent(eventID string) error
```

**Parameters:**
- `eventID`: ID of event to delete

**Returns:** `error` - Any error that occurred

---

## 💾 Database Package

### `pkg/database/database.go`

#### `Database`
File-based database for storing user interactions.

```go
type Database struct {
    dataDir     string
    interactions map[int64][]types.Interaction
    mu          sync.RWMutex
}
```

#### `NewDatabase()`
Creates a new database instance.

```go
func NewDatabase(dataDir string) (*Database, error)
```

**Parameters:**
- `dataDir`: Directory to store data files

**Returns:** `(*Database, error)` - Database instance and any error

**Storage:**
- Uses JSON files for persistence
- In-memory map for quick access
- Thread-safe with mutex protection

#### `AddInteraction()`
Stores a new user-AI interaction.

```go
func (d *Database) AddInteraction(userID int64, userMsg, aiResp, action string) error
```

**Parameters:**
- `userID`: Telegram user ID
- `userMsg`: User's message
- `aiResp`: AI's response
- `action`: Action taken

**Returns:** `error` - Any error that occurred

#### `GetUserContext()`
Retrieves recent conversation context for a user.

```go
func (d *Database) GetUserContext(userID int64, messageCount int) string
```

**Parameters:**
- `userID`: Telegram user ID
- `messageCount`: Number of recent messages to include

**Returns:** `string` - Formatted conversation context

**Context Format:**
```
Previous conversation:
- User: What's on my calendar today?
- Assistant: I'll check your calendar for today.
- User: Create a meeting tomorrow at 2pm
- Assistant: I'll create a meeting for tomorrow at 2pm.
```

#### `GetUserInteractions()`
Retrieves user interactions for analysis.

```go
func (d *Database) GetUserInteractions(userID int64, limit int) []types.Interaction
```

**Parameters:**
- `userID`: Telegram user ID
- `limit`: Maximum number of interactions to return

**Returns:** `[]types.Interaction` - User's interaction history

#### `GetUserStats()`
Retrieves user interaction statistics.

```go
func (d *Database) GetUserStats(userID int64) (int, time.Time, error)
```

**Parameters:**
- `userID`: Telegram user ID

**Returns:** `(int, time.Time, error)` - Interaction count, last interaction time, and any error

#### `Backup()`
Creates a backup of the database.

```go
func (d *Database) Backup() error
```

**Returns:** `error` - Any error that occurred

#### `Cleanup()`
Removes old interactions to prevent database bloat.

```go
func (d *Database) Cleanup(maxAge time.Duration) error
```

**Parameters:**
- `maxAge`: Maximum age of interactions to keep

**Returns:** `error` - Any error that occurred

---

## 📱 Telegram Package

### `pkg/telegram/bot.go`

#### `Bot`
Telegram bot wrapper.

```go
type Bot struct {
    bot *tgbotapi.BotAPI
}
```

#### `NewBot()`
Creates a new Telegram bot instance.

```go
func NewBot(token string) (*Bot, error)
```

**Parameters:**
- `token`: Telegram bot token

**Returns:** `(*Bot, error)` - Bot instance and any error

#### `SendMessage()`
Sends a message to a specific chat.

```go
func (t *Bot) SendMessage(chatID int64, text string) error
```

**Parameters:**
- `chatID`: Target chat ID
- `text`: Message text to send

**Returns:** `error` - Any error that occurred

**Features:**
- Automatically splits long messages (>4096 characters)
- Logs all message operations
- Returns message ID for tracking

#### `GetUpdatesChan()`
Gets the channel for receiving Telegram updates.

```go
func (t *Bot) GetUpdatesChan() tgbotapi.UpdatesChannel
```

**Returns:** `tgbotapi.UpdatesChannel` - Channel for receiving updates

#### `GetBotInfo()`
Retrieves bot information.

```go
func (t *Bot) GetBotInfo() tgbotapi.User
```

**Returns:** `tgbotapi.User` - Bot user information

#### `AnswerCallbackQuery()`
Answers callback queries from inline keyboards.

```go
func (t *Bot) AnswerCallbackQuery(callbackQueryID, text string) error
```

**Parameters:**
- `callbackQueryID`: ID of the callback query
- `text`: Text to show to the user

**Returns:** `error` - Any error that occurred

#### `EditMessageText()`
Edits an existing message.

```go
func (t *Bot) EditMessageText(chatID int64, messageID int, newText string) error
```

**Parameters:**
- `chatID`: Chat ID containing the message
- `messageID`: ID of the message to edit
- `newText`: New text content

**Returns:** `error` - Any error that occurred

#### `DeleteMessage()`
Deletes a message.

```go
func (t *Bot) DeleteMessage(chatID int64, messageID int) error
```

**Parameters:**
- `chatID`: Chat ID containing the message
- `messageID`: ID of the message to delete

**Returns:** `error` - Any error that occurred

#### `sendLongMessage()`
Internal method for sending long messages by splitting them.

```go
func (t *Bot) sendLongMessage(chatID int64, text string) error
```

**Parameters:**
- `chatID`: Target chat ID
- `text`: Long text to split and send

**Returns:** `error` - Any error that occurred

**Splitting Logic:**
- Splits on sentence boundaries when possible
- Ensures each part is under 4096 characters
- Sends parts sequentially

---

## ⚙️ Config Package

### `pkg/config/config.go`

#### `Config`
Application configuration structure.

```go
type Config struct {
    TelegramToken    string
    OpenAIKey        string
    GoogleCreds      string
    CalendarID       string
    Port             string
}
```

#### `Load()`
Loads configuration from environment variables.

```go
func Load() (*Config, error)
```

**Returns:** `(*Config, error)` - Configuration and any error

**Environment Variables:**
- `TELEGRAM_TOKEN`: Telegram bot token
- `OPENAI_API_KEY`: OpenAI API key
- `GOOGLE_CREDENTIALS_FILE`: Path to Google credentials JSON
- `GOOGLE_CALENDAR_ID`: Google Calendar ID
- `PORT`: HTTP server port (optional, defaults to 8080)

#### `MaskToken()`
Utility function to mask sensitive tokens in logs.

```go
func MaskToken(token string) string
```

**Parameters:**
- `token`: Token to mask

**Returns:** `string` - Masked token (shows first 4 and last 4 characters)

---

## 🚀 Main Application

### `cmd/bot/main.go`

#### `Bot`
Main bot structure that coordinates all components.

```go
type Bot struct {
    aiAgent     *ai.Agent
    telegramBot *telegram.Bot
    config      *config.Config
}
```

#### `NewBot()`
Creates a new bot instance with all components.

```go
func NewBot(cfg *config.Config) (*Bot, error)
```

**Parameters:**
- `cfg`: Application configuration

**Returns:** `(*Bot, error)` - Bot instance and any error

**Initialization Flow:**
1. Creates OpenAI service
2. Creates Google Calendar service
3. Creates Telegram bot
4. Creates database
5. Creates AI agent
6. Returns configured bot

#### `handleMessage()`
Processes incoming Telegram messages.

```go
func (b *Bot) handleMessage(update tgbotapi.Update)
```

**Parameters:**
- `update`: Telegram update containing the message

**Processing:**
- Extracts user ID, message text, and chat ID
- Delegates to AI agent for processing
- Logs all operations

#### `startBot()`
Starts the Telegram bot.

```go
func (b *Bot) startBot() error
```

**Returns:** `error` - Any error that occurred

**Operation:**
- Gets updates channel from Telegram bot
- Processes updates concurrently using goroutines
- Logs bot startup

#### `main()`
Application entry point.

**Flow:**
1. Loads configuration
2. Creates bot instance
3. Starts bot
4. Handles any fatal errors

---

## 🔧 Error Handling

### Error Types

#### `ConfigError`
Configuration-related errors.

```go
type ConfigError struct {
    Field string
    Err   error
}

func (e *ConfigError) Error() string
func (e *ConfigError) Unwrap() error
```

#### `CalendarError`
Calendar operation errors.

```go
type CalendarError struct {
    Operation string
    Err       error
}

func (e *CalendarError) Error() string
func (e *CalendarError) Unwrap() error
```

### Error Handling Patterns

#### 1. **Wrapping Errors**
```go
return fmt.Errorf("failed to create calendar service: %w", err)
```

#### 2. **Contextual Errors**
```go
return &CalendarError{
    Operation: "GetEvents",
    Err:       err,
}
```

#### 3. **Graceful Degradation**
```go
if err := b.aiAgent.ProcessUserMessage(userID, chatID, message); err != nil {
    log.Printf("Error processing message for user %d: %v", userID, err)
    // Send user-friendly error message
    errorMsg := "Sorry, I encountered an error processing your request. Please try again."
    if err := b.telegramBot.SendMessage(chatID, errorMsg); err != nil {
        log.Printf("Failed to send error message: %v", err)
    }
}
```

---

## 📊 Logging

### Log Format
```
2024/01/15 10:30:45 [INFO] Creating bot with config: Telegram=1234...abcd, OpenAI=sk-...abcd, GoogleCreds=./credentials/google-credentials.json, CalendarID=primary
2024/01/15 10:30:45 [INFO] OpenAI service created successfully
2024/01/15 10:30:45 [INFO] Google Calendar service created successfully
```

### Log Levels
- **INFO**: Normal operation information
- **ERROR**: Error conditions
- **DEBUG**: Detailed debugging information (when enabled)

### Context Tracking
All log entries include:
- User ID (when available)
- Chat ID (when available)
- Operation being performed
- Success/failure status

---

*Next: [Configuration](configuration.md) - Environment variables and settings*

