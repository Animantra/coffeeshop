#!/bin/bash

exec > /var/log/user-data.log 2>&1
set -x

sudo apt-get update
sudo apt-get upgrade -y
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

sudo apt-get install -y docker.io git curl
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker ubuntu

curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-amd64
sudo install -o root -g root -m 0755 ./kind /usr/local/bin/kind


cat <<EOF > /tmp/kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30080
    hostPort: 80
    listenAddress: "0.0.0.0"
  - containerPort: 32000
    hostPort: 3000
    listenAddress: "0.0.0.0"
  - containerPort: 32090
    hostPort: 9090
    listenAddress: "0.0.0.0"
  - containerPort: 31888
    hostPort: 8888
    listenAddress: "0.0.0.0"
EOF

sleep 5

kind create cluster --name coffeeshop-kube --config /tmp/kind-config.yaml

mkdir -p /home/ubuntu/.kube
kind get kubeconfig --name coffeeshop-kube > /home/ubuntu/.kube/config
chown -R ubuntu:ubuntu /home/ubuntu/.kube

cd /home/ubuntu
git clone https://github.com/Animantra/coffeeshop
chown -R ubuntu:ubuntu /home/ubuntu/coffeeshop
