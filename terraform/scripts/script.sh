#!/bin/bash

exec > /var/log/user-data.log 2>&1

sudo apt-get update
sudo apt-get upgrade -y
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

sudo apt-get install -y docker.io docker-compose-v2 git curl
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker ubuntu

cd /home/ubuntu
git clone https://github.com/Animantra/coffeeshop

sudo chown -R ubuntu:ubuntu /home/ubuntu/coffeeshop

cd /home/ubuntu/coffeeshop
sudo docker compose up -d