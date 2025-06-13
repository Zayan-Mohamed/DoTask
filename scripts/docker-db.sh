#!/bin/bash
# DoTask Database Management Script

DB_CONTAINER="dotask-postgres"
DB_USER="postgres"
DB_NAME="dotask"

echo "🗄️ DoTask Database Management"
echo "============================="

# Check if database container is running
if ! docker ps --filter name=$DB_CONTAINER --format "{{.Names}}" | grep -q $DB_CONTAINER; then
    echo "❌ Database container '$DB_CONTAINER' is not running."
    echo "Please start the DoTask environment first."
    exit 1
fi

case "$1" in
    "backup")
        echo "📦 Creating database backup..."
        BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
        docker exec $DB_CONTAINER pg_dump -U $DB_USER $DB_NAME > $BACKUP_FILE
        echo "✅ Backup created: $BACKUP_FILE"
        ;;
    
    "restore")
        if [ -z "$2" ]; then
            echo "❌ Please provide backup file: ./scripts/docker-db.sh restore <backup_file>"
            exit 1
        fi
        echo "🔄 Restoring database from $2..."
        docker exec -i $DB_CONTAINER psql -U $DB_USER $DB_NAME < $2
        echo "✅ Database restored from $2"
        ;;
    
    "reset")
        echo "⚠️ This will delete all data! Are you sure? (y/N)"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            echo "🗑️ Resetting database..."
            docker exec $DB_CONTAINER psql -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME;"
            docker exec $DB_CONTAINER psql -U $DB_USER -c "CREATE DATABASE $DB_NAME;"
            echo "✅ Database reset completed"
        else
            echo "❌ Database reset cancelled"
        fi
        ;;
    
    "shell"|"psql")
        echo "🔌 Connecting to database shell..."
        docker exec -it $DB_CONTAINER psql -U $DB_USER $DB_NAME
        ;;
    
    "logs")
        echo "📋 Database logs:"
        docker logs $DB_CONTAINER --tail=50
        ;;
    
    "status")
        echo "📊 Database Status:"
        docker exec $DB_CONTAINER pg_isready -U $DB_USER -d $DB_NAME
        echo ""
        echo "🔢 Database Size:"
        docker exec $DB_CONTAINER psql -U $DB_USER $DB_NAME -c "SELECT pg_size_pretty(pg_database_size('$DB_NAME')) as database_size;"
        echo ""
        echo "📊 Table Information:"
        docker exec $DB_CONTAINER psql -U $DB_USER $DB_NAME -c "\dt"
        ;;
    
    "migrate")
        echo "🔄 Running database migrations..."
        # Check if backend container is available for migrations
        if docker ps --filter name=dotask-backend --format "{{.Names}}" | grep -q dotask-backend; then
            docker exec dotask-backend sh -c "cd /root && ./main migrate"
            echo "✅ Migrations completed"
        else
            echo "❌ Backend container not found. Please ensure the application is running."
        fi
        ;;
    
    *)
        echo "Usage: $0 {backup|restore <file>|reset|shell|logs|status|migrate}"
        echo ""
        echo "Commands:"
        echo "  backup          - Create a database backup"
        echo "  restore <file>  - Restore database from backup file"
        echo "  reset           - Reset database (WARNING: deletes all data)"
        echo "  shell           - Open database shell (psql)"
        echo "  logs            - Show database logs"
        echo "  status          - Show database status and information"
        echo "  migrate         - Run database migrations"
        ;;
esac
