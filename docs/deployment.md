# Deployment

## 🚀 Deployment Options

The Calendar Assistant Bot can be deployed using various methods depending on your infrastructure and requirements. This guide covers Docker, Docker Compose, and production deployment strategies.

## 🐳 Docker Deployment

### Prerequisites

- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)

### Quick Start with Docker

#### 1. **Build the Image**
```bash
# Build the Docker image
docker build -t calendar-assistant-bot:latest .

# Verify the image was created
docker images | grep calendar-assistant-bot
```

#### 2. **Run the Container**
```bash
# Run with environment variables
docker run -d \
  --name calendar-bot \
  -e TELEGRAM_TOKEN=your_token_here \
  -e OPENAI_API_KEY=your_openai_key_here \
  -e GOOGLE_CREDENTIALS_FILE=/app/credentials/google-credentials.json \
  -e GOOGLE_CALENDAR_ID=primary \
  -v $(pwd)/credentials:/app/credentials:ro \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  calendar-assistant-bot:latest
```

#### 3. **Check Container Status**
```bash
# View running containers
docker ps

# Check container logs
docker logs calendar-bot

# Monitor resource usage
docker stats calendar-bot
```

### Docker Compose Deployment

#### 1. **Create docker-compose.yml**
```yaml
version: '3.8'

services:
  calendar-bot:
    build: .
    image: calendar-assistant-bot:latest
    container_name: calendar-assistant-bot
    restart: unless-stopped
    environment:
      - TELEGRAM_TOKEN=${TELEGRAM_TOKEN}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - GOOGLE_CREDENTIALS_FILE=/app/credentials/google-credentials.json
      - GOOGLE_CALENDAR_ID=${GOOGLE_CALENDAR_ID}
      - PORT=8080
    volumes:
      - ./credentials:/app/credentials:ro
      - ./data:/app/data
      - ./logs:/app/logs
    ports:
      - "8080:8080"
    networks:
      - bot-network

networks:
  bot-network:
    driver: bridge
```

#### 2. **Create .env File**
```bash
# .env
TELEGRAM_TOKEN=your_telegram_bot_token_here
OPENAI_API_KEY=your_openai_api_key_here
GOOGLE_CALENDAR_ID=primary
```

#### 3. **Deploy with Docker Compose**
```bash
# Start the service
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f calendar-bot

# Stop the service
docker-compose down
```

### Dockerfile Analysis

#### Multi-Stage Build
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/bot

# Production stage
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy credentials directory
COPY --from=builder /app/credentials ./credentials

# Create data and logs directories
RUN mkdir -p data logs

# Expose port (if implementing webhooks)
EXPOSE 8080

# Run the application
CMD ["./main"]
```

#### Build Optimizations
```dockerfile
# Use .dockerignore to exclude unnecessary files
# .dockerignore
.git
.gitignore
README.md
docs/
*.md
.env
.env.*
data/
logs/
credentials/
Dockerfile*
docker-compose*.yml
Makefile
debug.sh
```

## 🏭 Production Deployment

### Production Considerations

#### 1. **Security**
- Use secrets management (Docker secrets, Kubernetes secrets)
- Restrict container capabilities
- Use non-root user
- Implement proper logging and monitoring

#### 2. **Scalability**
- Use load balancers for multiple instances
- Implement health checks
- Use persistent storage for data
- Monitor resource usage

#### 3. **Reliability**
- Implement restart policies
- Use health checks
- Monitor application metrics
- Implement proper error handling

### Production Docker Compose

#### Enhanced docker-compose.prod.yml
```yaml
version: '3.8'

services:
  calendar-bot:
    build: .
    image: calendar-assistant-bot:latest
    container_name: calendar-assistant-bot-prod
    restart: unless-stopped
    environment:
      - TELEGRAM_TOKEN=${TELEGRAM_TOKEN}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - GOOGLE_CREDENTIALS_FILE=/app/credentials/google-credentials.json
      - GOOGLE_CALENDAR_ID=${GOOGLE_CALENDAR_ID}
      - PORT=8080
    volumes:
      - ./credentials:/app/credentials:ro
      - calendar-data:/app/data
      - calendar-logs:/app/logs
    ports:
      - "8080:8080"
    networks:
      - bot-network
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  nginx:
    image: nginx:alpine
    container_name: calendar-bot-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - calendar-bot
    networks:
      - bot-network

volumes:
  calendar-data:
    driver: local
  calendar-logs:
    driver: local

networks:
  bot-network:
    driver: bridge
```

### Nginx Configuration

#### Reverse Proxy Setup
```nginx
# nginx/nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream calendar-bot {
        server calendar-bot:8080;
    }

    server {
        listen 80;
        server_name your-domain.com;
        
        # Redirect HTTP to HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name your-domain.com;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;

        # Security headers
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";

        location / {
            proxy_pass http://calendar-bot;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /health {
            proxy_pass http://calendar-bot/health;
            access_log off;
        }
    }
}
```

## ☸️ Kubernetes Deployment

### Prerequisites

- **Kubernetes cluster**: Minikube, Docker Desktop, or cloud provider
- **kubectl**: [Install kubectl](https://kubernetes.io/docs/tasks/tools/)
- **Helm** (optional): [Install Helm](https://helm.sh/docs/intro/install/)

### Basic Kubernetes Deployment

#### 1. **Create Namespace**
```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: calendar-bot
```

#### 2. **Create ConfigMap**
```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: calendar-bot-config
  namespace: calendar-bot
data:
  GOOGLE_CREDENTIALS_FILE: /app/credentials/google-credentials.json
  GOOGLE_CALENDAR_ID: primary
  PORT: "8080"
```

#### 3. **Create Secrets**
```yaml
# k8s/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: calendar-bot-secrets
  namespace: calendar-bot
type: Opaque
data:
  telegram-token: <base64-encoded-token>
  openai-api-key: <base64-encoded-key>
```

#### 4. **Create Deployment**
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: calendar-bot
  namespace: calendar-bot
  labels:
    app: calendar-bot
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
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        env:
        - name: TELEGRAM_TOKEN
          valueFrom:
            secretKeyRef:
              name: calendar-bot-secrets
              key: telegram-token
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: calendar-bot-secrets
              key: openai-api-key
        - name: GOOGLE_CREDENTIALS_FILE
          valueFrom:
            configMapKeyRef:
              name: calendar-bot-config
              key: GOOGLE_CREDENTIALS_FILE
        - name: GOOGLE_CALENDAR_ID
          valueFrom:
            configMapKeyRef:
              name: calendar-bot-config
              key: GOOGLE_CALENDAR_ID
        - name: PORT
          valueFrom:
            configMapKeyRef:
              name: calendar-bot-config
              key: PORT
        volumeMounts:
        - name: google-credentials
          mountPath: /app/credentials
          readOnly: true
        - name: data-volume
          mountPath: /app/data
        - name: logs-volume
          mountPath: /app/logs
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: google-credentials
        secret:
          secretName: google-credentials
      - name: data-volume
        persistentVolumeClaim:
          claimName: calendar-bot-data-pvc
      - name: logs-volume
        persistentVolumeClaim:
          claimName: calendar-bot-logs-pvc
```

#### 5. **Create Service**
```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: calendar-bot-service
  namespace: calendar-bot
spec:
  selector:
    app: calendar-bot
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
```

#### 6. **Create Ingress**
```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: calendar-bot-ingress
  namespace: calendar-bot
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: calendar-bot-service
            port:
              number: 80
```

#### 7. **Create Persistent Volume Claims**
```yaml
# k8s/pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: calendar-bot-data-pvc
  namespace: calendar-bot
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: calendar-bot-logs-pvc
  namespace: calendar-bot
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

### Deploy to Kubernetes

#### 1. **Apply All Resources**
```bash
# Create namespace
kubectl apply -f k8s/namespace.yaml

# Apply all resources
kubectl apply -f k8s/

# Check deployment status
kubectl get all -n calendar-bot
```

#### 2. **Monitor Deployment**
```bash
# Check pod status
kubectl get pods -n calendar-bot

# View pod logs
kubectl logs -f deployment/calendar-bot -n calendar-bot

# Check service
kubectl get svc -n calendar-bot

# Check ingress
kubectl get ingress -n calendar-bot
```

## 🔧 Health Checks

### Health Check Endpoint

#### Implementation
```go
// Add to main.go or create a separate health package
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    
    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
        "version":   "1.0.0",
        "uptime":    time.Since(startTime).String(),
    }
    
    json.NewEncoder(w).Encode(health)
}

func startHealthServer(port string) {
    http.HandleFunc("/health", healthHandler)
    
    go func() {
        log.Printf("Starting health check server on port %s", port)
        if err := http.ListenAndServe(":"+port, nil); err != nil {
            log.Printf("Health server error: %v", err)
        }
    }()
}
```

#### Usage
```bash
# Test health endpoint
curl http://localhost:8080/health

# Expected response
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:45Z",
  "version": "1.0.0",
  "uptime": "2h15m30s"
}
```

## 📊 Monitoring and Logging

### Logging Configuration

#### Structured Logging
```go
// Implement structured logging with levels
type Logger struct {
    level string
}

func (l *Logger) Info(format string, args ...interface{}) {
    if l.level == "debug" || l.level == "info" {
        log.Printf("[INFO] "+format, args...)
    }
}

func (l *Logger) Error(format string, args ...interface{}) {
    log.Printf("[ERROR] "+format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
    if l.level == "debug" {
        log.Printf("[DEBUG] "+format, args...)
    }
}
```

#### Log Rotation
```bash
# Use logrotate for log management
# /etc/logrotate.d/calendar-bot
/path/to/logs/*.log {
    daily
    missingok
    rotate 7
    compress
    delaycompress
    notifempty
    create 644 root root
    postrotate
        systemctl reload calendar-bot
    endscript
}
```

### Metrics Collection

#### Prometheus Metrics
```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    messageCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "calendar_bot_messages_total",
            Help: "Total number of messages processed",
        },
        []string{"user_id", "action"},
    )
    
    processingDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "calendar_bot_processing_duration_seconds",
            Help:    "Time spent processing messages",
            Buckets: prometheus.DefBuckets,
        },
        []string{"action"},
    )
)

func init() {
    prometheus.MustRegister(messageCounter)
    prometheus.MustRegister(processingDuration)
}

// In your message processing
func (a *Agent) ProcessUserMessage(userID, chatID int64, message string) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        processingDuration.WithLabelValues("process").Observe(duration)
    }()
    
    // ... processing logic
    
    messageCounter.WithLabelValues(fmt.Sprintf("%d", userID), "processed").Inc()
    return nil
}
```

## 🔄 CI/CD Pipeline

### GitHub Actions

#### Workflow File
```yaml
# .github/workflows/deploy.yml
name: Deploy Calendar Bot

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Test
      run: go test -v ./...
    
    - name: Build
      run: go build -v ./cmd/bot

  build-and-deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Build Docker image
      run: docker build -t calendar-assistant-bot:${{ github.sha }} .
    
    - name: Deploy to production
      run: |
        # Deploy to your production environment
        # This could be Docker Compose, Kubernetes, etc.
```

### GitLab CI

#### Pipeline Configuration
```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

test:
  stage: test
  image: golang:1.21-alpine
  script:
    - go test -v ./...
    - go build -v ./cmd/bot

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t calendar-assistant-bot:$CI_COMMIT_SHA .
    - docker tag calendar-assistant-bot:$CI_COMMIT_SHA calendar-assistant-bot:latest

deploy:
  stage: deploy
  image: alpine:latest
  script:
    - apk add --no-cache openssh-client
    - eval $(ssh-agent -s)
    - echo "$SSH_PRIVATE_KEY" | tr -d '\r' | ssh-add -
    - ssh $DEPLOY_USER@$DEPLOY_HOST "cd /opt/calendar-bot && docker-compose pull && docker-compose up -d"
  only:
    - main
```

## 🚨 Troubleshooting

### Common Issues

#### 1. **Container Won't Start**
```bash
# Check container logs
docker logs calendar-bot

# Check environment variables
docker exec calendar-bot env | grep -E "(TELEGRAM|OPENAI|GOOGLE)"

# Verify credentials file
docker exec calendar-bot ls -la /app/credentials/
```

#### 2. **Permission Issues**
```bash
# Fix file permissions
chmod 600 credentials/google-credentials.json
chmod 700 credentials/
chmod 755 data/ logs/

# Check container user
docker exec calendar-bot whoami
docker exec calendar-bot ls -la /app/
```

#### 3. **Network Issues**
```bash
# Check network connectivity
docker exec calendar-bot ping google.com
docker exec calendar-bot curl -I https://api.openai.com

# Check DNS resolution
docker exec calendar-bot nslookup api.openai.com
```

#### 4. **Resource Issues**
```bash
# Check resource usage
docker stats calendar-bot

# Check disk space
docker exec calendar-bot df -h

# Check memory usage
docker exec calendar-bot free -h
```

### Debug Commands

#### Docker Debug
```bash
# Interactive shell in container
docker exec -it calendar-bot /bin/sh

# Check running processes
docker exec calendar-bot ps aux

# Check file system
docker exec calendar-bot find /app -type f -name "*.json"
```

#### Kubernetes Debug
```bash
# Get pod details
kubectl describe pod <pod-name> -n calendar-bot

# Check events
kubectl get events -n calendar-bot --sort-by='.lastTimestamp'

# Port forward for debugging
kubectl port-forward <pod-name> 8080:8080 -n calendar-bot
```

## 📈 Performance Tuning

### Container Optimization

#### Resource Limits
```yaml
# docker-compose.yml
services:
  calendar-bot:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

#### JVM-like Tuning for Go
```dockerfile
# Dockerfile
ENV GOMAXPROCS=4
ENV GOGC=100
ENV GOMEMLIMIT=512MiB
```

### Application Optimization

#### Connection Pooling
```go
// Implement connection pooling for external APIs
type APIClient struct {
    client *http.Client
    pool   chan *http.Client
}

func NewAPIClient() *APIClient {
    return &APIClient{
        client: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        },
    }
}
```

---

*This completes the deployment documentation. For more details, refer to the [Development Guide](development.md) or [API Reference](api-reference.md).*

