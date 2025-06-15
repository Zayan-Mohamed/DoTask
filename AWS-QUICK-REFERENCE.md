# 📚 DoTask AWS Deployment Quick Reference

## 🚀 Deployment Commands

### Initial Setup (Run once)

```bash
# 1. SSH into EC2
ssh -i your-key.pem ec2-user@YOUR_ELASTIC_IP

# 2. Run EC2 setup
curl -sSL https://raw.githubusercontent.com/yourusername/yourrepo/main/scripts/ec2-setup.sh | bash
newgrp docker

# 3. Clone repository
cd /opt/dotask
git clone https://github.com/yourusername/dotask.git .

# 4. Configure environment
cp .env.aws.example .env
nano .env  # Edit with your values

# 5. Setup database
chmod +x scripts/*.sh
./scripts/setup-database.sh

# 6. Deploy application
docker-compose -f docker-compose.aws.yml up -d --build
```

### Daily Operations

```bash
# Check application status
./scripts/health-check.sh

# View logs
docker-compose -f docker-compose.aws.yml logs -f

# Restart application
docker-compose -f docker-compose.aws.yml restart

# Update application
git pull
docker-compose -f docker-compose.aws.yml up -d --build

# Backup database
./scripts/backup-database.sh

# Monitor resources
htop
docker stats
```

### Troubleshooting

```bash
# Check individual container logs
docker logs dotask-backend
docker logs dotask-frontend

# Access container shell
docker exec -it dotask-backend sh
docker exec -it dotask-frontend sh

# Rebuild specific service
docker-compose -f docker-compose.aws.yml up -d --build backend

# Clean up and restart
docker-compose -f docker-compose.aws.yml down
docker system prune -f
docker-compose -f docker-compose.aws.yml up -d --build
```

## 🔍 Health Check URLs

- **Frontend**: http://YOUR_ELASTIC_IP
- **Backend API**: http://YOUR_ELASTIC_IP:8080/health
- **GraphQL Playground**: http://YOUR_ELASTIC_IP:8080/playground

## 📊 Monitoring

### AWS CloudWatch (Optional)

- Navigate to CloudWatch in AWS Console
- Check EC2 metrics for CPU, Memory, Network
- Set up alarms for high CPU usage

### Application Logs

```bash
# Real-time logs
docker-compose -f docker-compose.aws.yml logs -f

# Backend logs only
docker-compose -f docker-compose.aws.yml logs -f backend

# Frontend logs only
docker-compose -f docker-compose.aws.yml logs -f frontend
```

## 🔒 Security Checklist

- [ ] RDS in private subnet
- [ ] Security groups properly configured
- [ ] SSH key secure
- [ ] Strong JWT secret
- [ ] Strong RDS password
- [ ] SSL certificate (if using domain)
- [ ] Regular backups
- [ ] Monitor logs for suspicious activity

## 💰 Cost Monitoring

### Free Tier Limits

- **EC2**: 750 hours/month (t2.micro)
- **RDS**: 750 hours/month (db.t3.micro)
- **Storage**: 30GB EBS + 20GB RDS
- **Data Transfer**: 15GB outbound

### Monitor Usage

- Check AWS Billing Dashboard monthly
- Set up billing alerts
- Monitor CloudWatch metrics

## 🆘 Emergency Procedures

### Application Down

```bash
# Quick restart
docker-compose -f docker-compose.aws.yml restart

# Full rebuild
docker-compose -f docker-compose.aws.yml down
docker-compose -f docker-compose.aws.yml up -d --build
```

### Database Issues

```bash
# Check RDS status in AWS Console
# Test connection
docker run --rm postgres:15-alpine pg_isready -h "$RDS_ENDPOINT" -p 5432 -U postgres

# Restore from backup (if available)
./scripts/restore-database.sh backup-file.sql
```

### Out of Disk Space

```bash
# Clean Docker
docker system prune -af
docker volume prune -f

# Check disk usage
df -h
du -sh /*
```

## 📞 Support Resources

- **AWS Support**: https://console.aws.amazon.com/support/
- **AWS Free Tier**: https://aws.amazon.com/free/
- **Docker Documentation**: https://docs.docker.com/
- **PostgreSQL Documentation**: https://www.postgresql.org/docs/
