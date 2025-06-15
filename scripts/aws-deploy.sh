#!/bin/bash
# DoTask AWS Deployment Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_step() { echo -e "${BLUE}>>> $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }

# Check if environment variables are set
check_env_vars() {
    print_step "Checking environment variables..."
    
    if [ -z "$RDS_ENDPOINT" ]; then
        print_error "RDS_ENDPOINT not set"
        exit 1
    fi
    
    if [ -z "$RDS_PASSWORD" ]; then
        print_error "RDS_PASSWORD not set"
        exit 1
    fi
    
    if [ -z "$PUBLIC_IP" ]; then
        print_error "PUBLIC_IP not set"
        exit 1
    fi
    
    if [ -z "$JWT_SECRET" ]; then
        print_warning "JWT_SECRET not set, generating one..."
        export JWT_SECRET=$(openssl rand -hex 32)
        echo "Generated JWT_SECRET: $JWT_SECRET"
        echo "Save this for future deployments!"
    fi
    
    print_success "Environment variables validated"
}

# Install Docker and Docker Compose
install_docker() {
    print_step "Installing Docker..."
    
    # Update packages
    sudo yum update -y
    
    # Install Docker
    sudo yum install -y docker
    sudo systemctl start docker
    sudo systemctl enable docker
    sudo usermod -a -G docker ec2-user
    
    # Install Docker Compose
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.21.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
    
    print_success "Docker installed successfully"
    print_warning "Please log out and log back in for Docker permissions to take effect"
}

# Clone repository and deploy
deploy_application() {
    print_step "Deploying DoTask application..."
    
    # Clone repository (you'll need to upload your code or use git)
    if [ ! -d "DoTask" ]; then
        print_error "DoTask directory not found. Please upload your application code first."
        exit 1
    fi
    
    cd DoTask
    
    # Stop any existing containers
    docker-compose -f docker-compose.aws.yml down || true
    
    # Build and start containers
    docker-compose -f docker-compose.aws.yml up --build -d
    
    print_success "Application deployed successfully!"
}

# Check application health
check_health() {
    print_step "Checking application health..."
    
    # Wait for services to start
    sleep 30
    
    # Check frontend
    if curl -s http://localhost/health > /dev/null; then
        print_success "Frontend is healthy"
    else
        print_error "Frontend health check failed"
    fi
    
    # Check backend
    if curl -s http://localhost/api/query -d '{"query":"query { __schema { types { name } } }"}' -H "Content-Type: application/json" > /dev/null; then
        print_success "Backend API is healthy"
    else
        print_error "Backend API health check failed"
    fi
}

# Setup SSL certificate (Let's Encrypt)
setup_ssl() {
    print_step "Setting up SSL certificate..."
    
    # Install certbot
    sudo yum install -y certbot python3-certbot-nginx
    
    if [ -n "$DOMAIN_NAME" ]; then
        # Stop nginx temporarily
        docker stop dotask-frontend || true
        
        # Get certificate
        sudo certbot certonly --standalone -d $DOMAIN_NAME --non-interactive --agree-tos --email admin@$DOMAIN_NAME
        
        # Start nginx again
        docker start dotask-frontend || true
        
        print_success "SSL certificate obtained for $DOMAIN_NAME"
    else
        print_warning "DOMAIN_NAME not set, skipping SSL setup"
    fi
}

# Configure nginx for production
update_nginx_config() {
    print_step "Updating NGINX configuration for AWS..."
    
    # Update nginx config to include SSL and proper backend proxying
    cat > do-task/nginx.aws.conf << 'EOF'
# AWS Production NGINX configuration for DoTask
user nginx;
worker_processes auto;
worker_rlimit_nofile 65535;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logging format
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    # Basic settings
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    server_tokens off;
    client_max_body_size 16M;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1000;
    gzip_comp_level 6;
    gzip_types
        application/atom+xml
        application/javascript
        application/json
        application/rss+xml
        application/vnd.ms-fontobject
        application/x-font-ttf
        application/x-web-app-manifest+json
        application/xhtml+xml
        application/xml
        font/opentype
        image/svg+xml
        image/x-icon
        text/css
        text/plain
        text/x-component;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=login:10m rate=5r/m;

    server {
        listen 80;
        server_name _;
        root /usr/share/nginx/html;
        index index.html;
        
        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header Referrer-Policy "no-referrer-when-downgrade" always;

        # Health check
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }

        # GraphQL API proxy
        location /api/query {
            limit_req zone=api burst=20 nodelay;
            
            proxy_pass http://backend:8080/query;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            
            # WebSocket support
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            
            # Timeouts
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # Frontend static files
        location / {
            try_files $uri $uri/ /index.html;
            
            # Cache static assets
            location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
                expires 1y;
                add_header Cache-Control "public, immutable";
            }
        }
    }
}
EOF

    print_success "NGINX configuration updated for AWS"
}

# Main execution
case "$1" in
    "install")
        install_docker
        ;;
    "deploy")
        check_env_vars
        update_nginx_config
        deploy_application
        check_health
        ;;
    "ssl")
        setup_ssl
        ;;
    "health")
        check_health
        ;;
    *)
        echo "Usage: $0 {install|deploy|ssl|health}"
        echo ""
        echo "Commands:"
        echo "  install  - Install Docker and Docker Compose"
        echo "  deploy   - Deploy the application"
        echo "  ssl      - Setup SSL certificate"
        echo "  health   - Check application health"
        echo ""
        echo "Required environment variables for deploy:"
        echo "  RDS_ENDPOINT - RDS database endpoint"
        echo "  RDS_PASSWORD - RDS database password"
        echo "  PUBLIC_IP    - EC2 public IP address"
        echo "  JWT_SECRET   - JWT secret key (optional, will generate if not set)"
        echo "  DOMAIN_NAME  - Domain name for SSL (optional)"
        ;;
esac
