#!/bin/bash
# DoTask AWS Quick Setup Script

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

print_header() {
    echo -e "${BOLD}${BLUE}"
    echo "╔═══════════════════════════════════════════════════════════════════╗"
    echo "║                     DoTask AWS Deployment                        ║"
    echo "║                    Quick Setup Assistant                         ║"
    echo "╚═══════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_step() { echo -e "${BOLD}${YELLOW}📋 $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }

# Check if we're on EC2
check_environment() {
    print_step "Checking environment..."
    
    if curl -s --max-time 2 http://169.254.169.254/latest/meta-data/instance-id > /dev/null 2>&1; then
        print_success "Running on AWS EC2 instance"
        return 0
    else
        print_info "Not running on EC2 - this script is designed for EC2 instances"
        echo "If you're setting up locally, please follow the manual AWS deployment guide."
        return 1
    fi
}

# Get EC2 metadata
get_ec2_info() {
    print_step "Gathering EC2 instance information..."
    
    INSTANCE_ID=$(curl -s http://169.254.169.254/latest/meta-data/instance-id)
    REGION=$(curl -s http://169.254.169.254/latest/meta-data/placement/region)
    PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)
    
    echo "Instance ID: $INSTANCE_ID"
    echo "Region: $REGION"
    echo "Public IP: $PUBLIC_IP"
    
    export PUBLIC_IP
}

# Setup environment variables
setup_environment() {
    print_step "Setting up environment variables..."
    
    # Check if .env.aws exists
    if [ -f ".env.aws" ]; then
        print_info "Found .env.aws file, sourcing it..."
        source .env.aws
    else
        print_info "Creating .env.aws from template..."
        if [ -f ".env.aws.template" ]; then
            cp .env.aws.template .env.aws
            print_error "Please edit .env.aws with your actual values and run this script again"
            print_info "Required values: RDS_ENDPOINT, RDS_PASSWORD"
            exit 1
        else
            print_error ".env.aws.template not found. Please check your DoTask directory."
            exit 1
        fi
    fi
    
    # Validate required variables
    if [ -z "$RDS_ENDPOINT" ] || [ -z "$RDS_PASSWORD" ]; then
        print_error "Missing required environment variables: RDS_ENDPOINT, RDS_PASSWORD"
        print_info "Please edit .env.aws with your actual values"
        exit 1
    fi
    
    # Generate JWT_SECRET if not set
    if [ -z "$JWT_SECRET" ]; then
        print_info "Generating JWT_SECRET..."
        JWT_SECRET=$(openssl rand -hex 32)
        echo "JWT_SECRET=$JWT_SECRET" >> .env.aws
        export JWT_SECRET
        print_success "Generated JWT_SECRET: $JWT_SECRET"
        print_info "This has been saved to .env.aws for future use"
    fi
    
    print_success "Environment variables configured"
}

# Install prerequisites
install_prerequisites() {
    print_step "Installing prerequisites..."
    
    # Update system
    sudo yum update -y
    
    # Install Docker
    if ! command -v docker &> /dev/null; then
        print_info "Installing Docker..."
        sudo yum install -y docker
        sudo systemctl start docker
        sudo systemctl enable docker
        sudo usermod -a -G docker ec2-user
        print_success "Docker installed"
    else
        print_success "Docker already installed"
    fi
    
    # Install Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        print_info "Installing Docker Compose..."
        sudo curl -L "https://github.com/docker/compose/releases/download/v2.21.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
        sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
        print_success "Docker Compose installed"
    else
        print_success "Docker Compose already installed"
    fi
    
    # Install other utilities
    sudo yum install -y git curl wget unzip
    
    print_success "Prerequisites installed"
}

# Deploy application
deploy_application() {
    print_step "Deploying DoTask application..."
    
    # Stop any existing containers
    if [ -f "docker-compose.aws.yml" ]; then
        docker-compose -f docker-compose.aws.yml down || true
    fi
    
    # Build and start containers
    docker-compose -f docker-compose.aws.yml up --build -d
    
    print_success "Application deployment started"
}

# Wait for services and test
test_deployment() {
    print_step "Testing deployment..."
    
    print_info "Waiting for services to start (this may take a few minutes)..."
    sleep 60
    
    # Test frontend
    if curl -s http://localhost/health > /dev/null; then
        print_success "Frontend is responding"
    else
        print_error "Frontend health check failed"
        docker-compose -f docker-compose.aws.yml logs frontend
    fi
    
    # Test backend
    local graphql_test=$(curl -s -X POST http://localhost/api/query \
        -H "Content-Type: application/json" \
        -d '{"query":"query { __schema { types { name } } }"}' | grep -o "types" || echo "")
    
    if [ "$graphql_test" = "types" ]; then
        print_success "Backend API is responding"
    else
        print_error "Backend API test failed"
        docker-compose -f docker-compose.aws.yml logs backend
    fi
}

# Show final information
show_completion_info() {
    print_header
    print_success "DoTask deployment completed!"
    echo ""
    echo -e "${BOLD}🌐 Access URLs:${NC}"
    echo "  • Application: http://$PUBLIC_IP"
    echo "  • Health Check: http://$PUBLIC_IP/health"
    echo "  • API Endpoint: http://$PUBLIC_IP/api/query"
    echo ""
    echo -e "${BOLD}📋 Management Commands:${NC}"
    echo "  • View logs: docker-compose -f docker-compose.aws.yml logs -f"
    echo "  • Restart: docker-compose -f docker-compose.aws.yml restart"
    echo "  • Stop: docker-compose -f docker-compose.aws.yml down"
    echo "  • Update: git pull && docker-compose -f docker-compose.aws.yml up --build -d"
    echo ""
    echo -e "${BOLD}🔐 Security:${NC}"
    echo "  • Your JWT_SECRET: $JWT_SECRET"
    echo "  • Database: $RDS_ENDPOINT"
    echo ""
    echo -e "${BOLD}📚 Documentation:${NC}"
    echo "  • Full guide: ./AWS-DEPLOYMENT.md"
    echo "  • Set up SSL: ./scripts/aws-deploy.sh ssl"
    echo ""
    if [ -n "$DOMAIN_NAME" ]; then
        echo -e "${YELLOW}🌍 Don't forget to point your domain ($DOMAIN_NAME) to $PUBLIC_IP${NC}"
        echo ""
    fi
}

# Main execution
main() {
    print_header
    
    if ! check_environment; then
        exit 1
    fi
    
    get_ec2_info
    setup_environment
    install_prerequisites
    
    # Check if user needs to log out for Docker permissions
    if ! docker ps &> /dev/null; then
        print_error "Docker permissions not yet active"
        print_info "Please log out and log back in, then run this script again"
        print_info "Command: exit (then SSH back in)"
        exit 1
    fi
    
    deploy_application
    test_deployment
    show_completion_info
}

# Help function
show_help() {
    echo "DoTask AWS Quick Setup"
    echo "====================="
    echo ""
    echo "This script will automatically set up DoTask on your AWS EC2 instance."
    echo ""
    echo "Prerequisites:"
    echo "  • Running on AWS EC2 instance"
    echo "  • RDS PostgreSQL database created"
    echo "  • .env.aws file with RDS_ENDPOINT and RDS_PASSWORD"
    echo ""
    echo "Usage:"
    echo "  $0           # Run full setup"
    echo "  $0 --help    # Show this help"
    echo ""
    echo "For detailed setup instructions, see: ./AWS-DEPLOYMENT.md"
}

# Script execution
case "$1" in
    "--help"|"-h"|"help")
        show_help
        ;;
    *)
        main
        ;;
esac
