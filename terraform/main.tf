provider "aws" {
  region     = "eu-north-1"
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}

resource "aws_security_group" "coffeeshop_sg" {
  name        = "coffeeshop-sg-new"
  description = "Allow web and monitoring traffic"

  dynamic "ingress" {
    for_each = [22, 80, 3000, 9090, 8888, 5001]
    content {
      from_port   = ingress.value
      to_port     = ingress.value
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "app_server" {
  ami           = var.ami_id
  instance_type = var.instance_type
  key_name      = var.key_name
  vpc_security_group_ids = [aws_security_group.coffeeshop_sg.id]

  root_block_device {
    volume_size = 20  
    volume_type = "gp3" 
  }
  user_data = <<-EOF
              #!/bin/bash
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
              git clone https://${var.github_token}@github.com/Animantra/coffeeshop.git
              
              sudo chown -R ubuntu:ubuntu /home/ubuntu/coffeeshop

              # cd /home/ubuntu/coffeeshop
              # sudo docker compose up -d
              EOF

  tags = {
    Name = "CoffeeShop-SRE-Server"
  }
}