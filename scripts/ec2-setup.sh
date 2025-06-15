#!/bin/bash
# EC2 Initial Setup Script for DoTask Deployment

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

print_step "Starting EC2 setup for DoTask deployment..."

# Update system
print_step "Updating system packages..."
sudo dnf update -y

# Install Docker
print_step "Installing Docker..."
sudo dnf install -y docker
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user

# Install Docker Compose
print_step "Installing Docker Compose..."
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Git
print_step "Installing Git..."
sudo dnf install -y git

# Install other utilities
print_step "Installing additional utilities..."
sudo dnf install -y htop curl wget unzip

# Create application directory
print_step "Creating application directory..."
sudo mkdir -p /opt/dotask
sudo chown ec2-user:ec2-user /opt/dotask

print_success "EC2 setup completed!"
print_step "Please log out and log back in for Docker group membership to take effect."
print_step "Then run: newgrp docker"
