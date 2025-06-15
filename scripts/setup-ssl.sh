#!/bin/bash
# SSL Setup Script for DoTask using Let's Encrypt

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_step() { echo -e "${BLUE}>>> $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

# Check if domain is provided
if [ -z "$1" ]; then
    print_error "Usage: $0 <domain-name> <email>"
    print_error "Example: $0 yourdomain.com your-email@domain.com"
    exit 1
fi

if [ -z "$2" ]; then
    print_error "Usage: $0 <domain-name> <email>"
    print_error "Example: $0 yourdomain.com your-email@domain.com"
    exit 1
fi

DOMAIN=$1
EMAIL=$2

print_step "Setting up SSL for domain: $DOMAIN"

# Install Certbot
print_step "Installing Certbot..."
sudo dnf install -y certbot python3-certbot-nginx

# Stop NGINX temporarily
print_step "Stopping NGINX for certificate generation..."
docker-compose -f docker-compose.aws.yml stop frontend

# Generate SSL certificate
print_step "Generating SSL certificate..."
sudo certbot certonly --standalone \
    --preferred-challenges http \
    --email "$EMAIL" \
    --agree-tos \
    --no-eff-email \
    -d "$DOMAIN"

# Update environment variables
print_step "Updating environment variables..."
sed -i "s/ENABLE_SSL=false/ENABLE_SSL=true/" .env
sed -i "s/DOMAIN_NAME=.*/DOMAIN_NAME=$DOMAIN/" .env

# Restart with SSL
print_step "Restarting application with SSL..."
docker-compose -f docker-compose.aws.yml up -d

# Set up auto-renewal
print_step "Setting up SSL certificate auto-renewal..."
echo "0 12 * * * /usr/bin/certbot renew --quiet && docker-compose -f /opt/dotask/docker-compose.aws.yml restart frontend" | sudo crontab -

print_success "SSL setup completed!"
print_success "Your application is now available at: https://$DOMAIN"
