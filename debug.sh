#!/bin/bash

# Debug script for Go Telegram AI Bot
# Usage: ./debug.sh [command]

CONTAINER_NAME="calendar-assistant-bot"

case "${1:-logs}" in
    "logs")
        echo "=== Viewing container logs ==="
        docker logs -f $CONTAINER_NAME
        ;;
    "logs-tail")
        echo "=== Viewing last 50 lines of logs ==="
        docker logs --tail 50 $CONTAINER_NAME
        ;;
    "exec")
        echo "=== Entering container shell ==="
        docker exec -it $CONTAINER_NAME /bin/sh
        ;;
    "env")
        echo "=== Container environment variables ==="
        docker exec $CONTAINER_NAME env | grep -E "(TELEGRAM|OPENAI|GOOGLE|PORT)"
        ;;
    "files")
        echo "=== Checking container files ==="
        docker exec $CONTAINER_NAME ls -la /app/
        docker exec $CONTAINER_NAME ls -la /app/credentials/ 2>/dev/null || echo "Credentials directory not found"
        docker exec $CONTAINER_NAME ls -la /app/logs/ 2>/dev/null || echo "Logs directory not found"
        ;;
    "status")
        echo "=== Container status ==="
        docker ps | grep $CONTAINER_NAME
        ;;
    "restart")
        echo "=== Restarting container ==="
        docker restart $CONTAINER_NAME
        ;;
    "rebuild")
        echo "=== Rebuilding and restarting container ==="
        docker-compose -f docker-compose.go.yml down
        docker-compose -f docker-compose.go.yml build --no-cache
        docker-compose -f docker-compose.go.yml up -d
        ;;
    "help")
        echo "Available commands:"
        echo "  logs      - View logs in real-time (default)"
        echo "  logs-tail - View last 50 lines of logs"
        echo "  exec      - Enter container shell"
        echo "  env       - Show environment variables"
        echo "  files     - Check container files"
        echo "  status    - Show container status"
        echo "  restart   - Restart container"
        echo "  rebuild   - Rebuild and restart container"
        echo "  help      - Show this help"
        ;;
    *)
        echo "Unknown command: $1"
        echo "Use './debug.sh help' for available commands"
        ;;
esac
