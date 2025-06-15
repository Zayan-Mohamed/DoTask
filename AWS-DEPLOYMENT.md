# 🚀 AWS Deployment Guide for DoTask

This guide will help you deploy the DoTask application to AWS using the free tier.

## 📋 Prerequisites

- AWS Account with free tier eligibility
- Domain name (optional, for SSL/HTTPS)
- Local computer with SSH client

## 🏗️ Architecture Overview

```
Internet → Elastic IP → EC2 (Docker) → RDS PostgreSQL
                       ├── Frontend (SvelteKit + NGINX)
                       └── Backend (Go GraphQL API)
```

## 💰 Cost Estimate (Free Tier)

- **EC2 t2.micro**: Free for 12 months (750 hours/month)
- **RDS t3.micro**: Free for 12 months (750 hours/month)
- **VPC, Security Groups**: Always free
- **Elastic IP**: Free when attached to running instance
- **Storage**: 30GB EBS, 20GB RDS (free tier limits)

## 🚀 Step-by-Step Deployment

### Step 1: Create AWS Infrastructure

#### 1.1 Create VPC

```bash
AWS Console → VPC → Create VPC
- Name: dotask-vpc
- IPv4 CIDR: 10.0.0.0/16
```

#### 1.2 Create Subnets

```bash
Public Subnet:
- Name: dotask-public-subnet
- AZ: us-east-1a
- CIDR: 10.0.1.0/24

Private Subnet 1:
- Name: dotask-private-subnet-1
- AZ: us-east-1a
- CIDR: 10.0.2.0/24

Private Subnet 2:
- Name: dotask-private-subnet-2
- AZ: us-east-1b
- CIDR: 10.0.3.0/24
```

#### 1.3 Create Internet Gateway

```bash
- Name: dotask-igw
- Attach to: dotask-vpc
```

#### 1.4 Create Route Table

```bash
- Name: dotask-public-rt
- VPC: dotask-vpc
- Routes: 0.0.0.0/0 → dotask-igw
- Associate: dotask-public-subnet
```

#### 1.5 Create Security Groups

```bash
EC2 Security Group (dotask-ec2-sg):
- SSH (22): Your IP
- HTTP (80): 0.0.0.0/0
- HTTPS (443): 0.0.0.0/0
- Custom (8080): 0.0.0.0/0

RDS Security Group (dotask-rds-sg):
- PostgreSQL (5432): dotask-ec2-sg
```

### Step 2: Create RDS Database

```bash
AWS Console → RDS → Create Database

Engine: PostgreSQL 15.x
Template: Free tier
Instance: db.t3.micro
Storage: 20GB gp2

Settings:
- DB identifier: dotask-database
- Username: postgres
- Password: [YourSecurePassword]
- Database name: dotask

Connectivity:
- VPC: dotask-vpc
- Subnet group: Create new (dotask-db-subnet-group)
- Security group: dotask-rds-sg
- Public access: No
```

**Important**: Save the RDS endpoint URL, you'll need it later!

### Step 3: Launch EC2 Instance

```bash
AWS Console → EC2 → Launch Instance

Name: dotask-server
AMI: Amazon Linux 2023
Instance type: t2.micro
Key pair: Create new or use existing

Network:
- VPC: dotask-vpc
- Subnet: dotask-public-subnet
- Auto-assign IP: Enable
- Security group: dotask-ec2-sg

Storage: 8GB gp2
```

### Step 4: Allocate Elastic IP

```bash
AWS Console → EC2 → Elastic IPs
1. Allocate Elastic IP
2. Associate with dotask-server instance
```

### Step 5: Connect and Setup Server

#### 5.1 SSH to your instance

```bash
ssh -i your-key.pem ec2-user@YOUR_ELASTIC_IP
```

#### 5.2 Install Docker

```bash
# Run the installation script
./aws-deploy.sh install
```

Or manually:

```bash
sudo yum update -y
sudo yum install -y docker git
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.21.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

**Important**: Log out and log back in after installing Docker!

### Step 6: Deploy Application

#### 6.1 Upload your code

You can either:

1. **Use git** (if your repo is public):

   ```bash
   git clone https://github.com/your-username/DoTask.git
   ```

2. **Use SCP to upload** (for private repos):
   ```bash
   # From your local machine
   scp -i your-key.pem -r /path/to/DoTask ec2-user@YOUR_ELASTIC_IP:~/
   ```

#### 6.2 Set environment variables

```bash
# SSH into your EC2 instance
export RDS_ENDPOINT="your-rds-endpoint.region.rds.amazonaws.com"
export RDS_PASSWORD="YourSecurePassword"
export PUBLIC_IP="YOUR_ELASTIC_IP"
export JWT_SECRET="your-jwt-secret-key"  # Optional, will generate if not set
export DOMAIN_NAME="your-domain.com"     # Optional, for SSL
```

#### 6.3 Deploy the application

```bash
cd DoTask
chmod +x scripts/aws-deploy.sh
./scripts/aws-deploy.sh deploy
```

### Step 7: Verify Deployment

#### 7.1 Check application health

```bash
./scripts/aws-deploy.sh health
```

#### 7.2 Test the application

```bash
# Frontend
curl http://YOUR_ELASTIC_IP/health

# Backend API
curl -X POST http://YOUR_ELASTIC_IP/api/query \
  -H "Content-Type: application/json" \
  -d '{"query":"query { __schema { types { name } } }"}'
```

#### 7.3 View in browser

Navigate to: `http://YOUR_ELASTIC_IP`

### Step 8: Setup SSL (Optional)

If you have a domain name:

```bash
# Point your domain to the Elastic IP in your DNS settings
# Then run:
export DOMAIN_NAME="your-domain.com"
./scripts/aws-deploy.sh ssl
```

## 🔧 Management Commands

### View logs

```bash
docker-compose -f docker-compose.aws.yml logs -f
```

### Restart services

```bash
docker-compose -f docker-compose.aws.yml restart
```

### Update application

```bash
git pull  # or upload new code
docker-compose -f docker-compose.aws.yml up --build -d
```

### Backup database

```bash
# Create backup
docker exec dotask-backend sh -c "PGPASSWORD='$RDS_PASSWORD' pg_dump -h $RDS_ENDPOINT -U postgres dotask" > backup.sql

# Upload to S3 (optional)
aws s3 cp backup.sql s3://your-backup-bucket/
```

## 📊 Monitoring and Maintenance

### Check resource usage

```bash
# Disk space
df -h

# Memory usage
free -h

# Docker containers
docker stats
```

### Log rotation

```bash
# Setup logrotate for Docker logs
sudo tee /etc/logrotate.d/docker > /dev/null <<EOF
/var/lib/docker/containers/*/*.log {
    rotate 7
    daily
    compress
    size 10M
    missingok
    delaycompress
    copytruncate
}
EOF
```

## 🚨 Troubleshooting

### Common Issues

1. **Connection refused to RDS**:

   - Check security group allows connection from EC2
   - Verify RDS endpoint is correct
   - Ensure RDS is in same VPC

2. **Application not accessible**:

   - Check security group allows HTTP/HTTPS
   - Verify Elastic IP is associated
   - Check Docker containers are running: `docker ps`

3. **Out of memory**:
   - Free tier has limited memory
   - Consider using swap file:
     ```bash
     sudo dd if=/dev/zero of=/swapfile bs=128M count=8
     sudo chmod 600 /swapfile
     sudo mkswap /swapfile
     sudo swapon /swapfile
     ```

### Logs

```bash
# Application logs
docker-compose -f docker-compose.aws.yml logs

# System logs
sudo journalctl -u docker

# Nginx logs
docker exec dotask-frontend tail -f /var/log/nginx/access.log
```

## 🔐 Security Best Practices

1. **Update regularly**:

   ```bash
   sudo yum update -y
   ```

2. **Limit SSH access**:

   - Use key-based authentication only
   - Restrict SSH to your IP in security group

3. **Use HTTPS**:

   - Set up SSL certificate
   - Redirect HTTP to HTTPS

4. **Database security**:
   - Use strong passwords
   - Keep database in private subnet
   - Enable automated backups

## 💵 Cost Optimization

1. **Stop instances when not needed**:

   - Free tier gives 750 hours/month
   - Stop development instances overnight

2. **Monitor usage**:

   - Set up billing alerts
   - Check AWS Free Tier usage dashboard

3. **Clean up resources**:
   - Delete unattached EBS volumes
   - Remove unused snapshots
   - Clean up Docker images: `docker system prune`

## 🎯 Next Steps

1. **Set up CI/CD**: Use GitHub Actions with AWS CodeDeploy
2. **Add monitoring**: CloudWatch, Application Load Balancer
3. **Scale**: Move to ECS or EKS for production
4. **Backup strategy**: Automated RDS snapshots, S3 backups

## 📞 Support

If you encounter issues:

1. Check the troubleshooting section
2. Review AWS documentation
3. Check application logs
4. Verify security group configurations

---

**Remember**: This setup is for development/testing. For production, consider using:

- Application Load Balancer
- Auto Scaling Groups
- RDS Multi-AZ
- CloudFront CDN
- WAF for security
