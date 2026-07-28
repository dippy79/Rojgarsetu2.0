#!/bin/bash
# VPS Setup Script for RojgarSetu - Ubuntu 22.04

set -e

echo "🔧 Setting up RojgarSetu VPS (Ubuntu 22.04)..."

# Update system
echo "📦 Updating system packages..."
apt update &amp;&amp; apt upgrade -y

# Install Docker
echo "🐳 Installing Docker..."
curl -fsSL https://get.docker.com | sh
systemctl enable docker
systemctl start docker

# Docker Compose v2 plugin
echo "🔌 Installing Docker Compose v2..."
apt install docker-compose-plugin -y

# Nginx + Certbot
echo "🌐 Installing Nginx + Certbot..."
apt install nginx certbot python3-certbot-nginx -y
systemctl enable nginx
systemctl start nginx

# Git
echo "📂 Installing Git..."
apt install git -y

# UFW firewall
echo "🛡️ Configuring UFW firewall..."
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

echo "✅ VPS setup complete!"
echo "Next steps:"
echo "1. git clone your-repo /opt/rojgarsetu"
echo "2. cd /opt/rojgarsetu &amp;&amp; cp .env.production.example .env.production"
echo "3. Edit .env.production with real secrets"
echo "4. ./deploy.sh"
echo "5. certbot --nginx -d api.yourdomain.com"

