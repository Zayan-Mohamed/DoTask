# DoTask Docker Setup

This directory contains the complete Docker configuration for the DoTask application using Context 7 architecture.

## Quick Start

### Production Environment

```bash
# Start production environment
./scripts/docker-prod.sh

# Access the application
# Frontend: http://localhost
# API: http://localhost/api
# Health Check: http://localhost/health
```

### Development Environment

```bash
# Start development environment with hot reload
./scripts/docker-dev.sh

# Access the application
# Frontend: http://localhost:5173
# Backend: http://localhost:8080
# GraphQL Playground: http://localhost:8080
# Database: localhost:5432
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   NGINX Proxy   │    │   SvelteKit     │    │   Go Backend    │
│   (Port 80)     │◄──►│   Frontend      │◄──►│   (GraphQL)     │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                               ┌─────────────────┐
                                               │   PostgreSQL    │
                                               │   Database      │
                                               │                 │
                                               └─────────────────┘
```

## Services

### 1. NGINX Reverse Proxy

- **Port**: 80 (production)
- **Features**:
  - Static file serving
  - API proxying with `/api` prefix
  - Rate limiting
  - Security headers
  - GZIP compression
  - Health checks

### 2. SvelteKit Frontend

- **Development**: Port 5173 (direct access)
- **Production**: Served through NGINX
- **Features**:
  - SSR/SPA hybrid
  - Hot reload in development
  - Optimized production builds

### 3. Go Backend (GraphQL API)

- **Port**: 8080 (internal), `/api` (via NGINX)
- **Features**:
  - GraphQL API with Playground
  - JWT authentication
  - Database migrations
  - Hot reload in development

### 4. PostgreSQL Database

- **Port**: 5432
- **Features**:
  - Persistent data storage
  - Health checks
  - Automatic initialization

## Management Scripts

### Environment Management

```bash
# Start development environment
./scripts/docker-dev.sh

# Start production environment
./scripts/docker-prod.sh

# Check status and logs
./scripts/docker-status.sh
./scripts/docker-status.sh --logs
./scripts/docker-status.sh --follow

# Cleanup containers
./scripts/docker-cleanup.sh          # Keep data
./scripts/docker-cleanup.sh --full   # Remove all data
```

### Database Management

```bash
# Database operations
./scripts/docker-db.sh status        # Show database info
./scripts/docker-db.sh shell         # Open database shell
./scripts/docker-db.sh backup        # Create backup
./scripts/docker-db.sh restore file  # Restore from backup
./scripts/docker-db.sh reset         # Reset database
./scripts/docker-db.sh migrate       # Run migrations
./scripts/docker-db.sh logs          # Show database logs
```

## Configuration Files

### Docker Compose

- `docker-compose.yml` - Production environment
- `docker-compose.dev.yml` - Development environment with hot reload

### Dockerfiles

- `do-task/Dockerfile` - Production frontend build with NGINX
- `do-task/Dockerfile.dev` - Development frontend with hot reload
- `do-task-backend/Dockerfile` - Production backend build
- `do-task-backend/Dockerfile.dev` - Development backend with Air hot reload

### NGINX Configuration

- `do-task/nginx.conf` - Reverse proxy configuration with security and performance optimizations

### Environment Files

- `.env.docker` - Docker-specific environment variables
- Backend uses embedded environment variables in Dockerfiles

## Health Checks

All services include health checks:

- **Frontend**: HTTP check on `/health`
- **Backend**: HTTP check on root endpoint
- **Database**: PostgreSQL ready check

## Security Features

### NGINX Security Headers

- X-Frame-Options: SAMEORIGIN
- X-XSS-Protection: 1; mode=block
- X-Content-Type-Options: nosniff
- Content Security Policy
- Referrer Policy

### Rate Limiting

- API endpoints: 10 requests/second
- Login endpoints: 5 requests/minute
- Burst allowance configured

### CORS Configuration

- Proper CORS headers for frontend-backend communication
- Configurable origins

## Development Features

### Hot Reload

- **Frontend**: Vite dev server with HMR
- **Backend**: Air for Go hot reload
- **Volume mounts**: Source code changes reflected immediately

### Debug Mode

- Detailed logging in development
- GraphQL Playground access
- Source maps enabled

## Production Optimizations

### Multi-stage Builds

- Minimal final images
- Build dependencies excluded from runtime
- Alpine Linux base images

### Static Asset Optimization

- GZIP compression
- Long-term caching headers
- Optimized asset delivery

### Performance

- Keep-alive connections
- Connection pooling
- Optimized worker processes

## Networking

### Production Network

- `dotask-network`: Bridge network for service communication
- Services communicate by container name
- External access only through NGINX

### Development Network

- `dotask-dev-network`: Development bridge network
- Direct port access for debugging
- Same internal communication patterns

## Data Persistence

### Volumes

- `postgres_data`: Production database storage
- `postgres_dev_data`: Development database storage
- Automatic volume management

### Backups

- Manual backup script provided
- Timestamped backup files
- Easy restore functionality

## Troubleshooting

### Common Issues

1. **Port Conflicts**: Ensure ports 80, 5173, 8080, 5432 are available
2. **Database Connection**: Wait for health checks to pass
3. **CORS Errors**: Check CORS_ORIGINS environment variable
4. **Build Failures**: Check Dockerfile syntax and dependencies

### Debugging Commands

```bash
# Check container logs
docker logs dotask-frontend
docker logs dotask-backend
docker logs dotask-postgres

# Inspect container
docker exec -it dotask-backend sh

# Check networks
docker network ls
docker network inspect dotask-network

# Check volumes
docker volume ls
docker volume inspect dotask_postgres_data
```

## Context 7 Features

This implementation follows Context 7 principles:

1. **Intelligent Orchestration**: Smart health checks and dependency management
2. **Security First**: Multiple layers of security controls
3. **Performance Optimized**: Caching, compression, and optimized builds
4. **Developer Experience**: Hot reload, easy scripts, comprehensive logging
5. **Production Ready**: Multi-stage builds, proper networking, monitoring
6. **Maintainable**: Clear structure, documentation, management scripts
7. **Scalable Foundation**: Ready for load balancing and horizontal scaling
