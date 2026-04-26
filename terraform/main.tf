resource "aws_instance" "app_server" {
    ami           = "ami-0014ce3e52359afbd" 
  instance_type = "c7i-flex.large"            
    key_name      = "coffeeshop-key"       

    vpc_security_group_ids = [aws_security_group.app_sg.id]

    tags = { Name = "Coffeeshop-Server" }
}

resource "aws_security_group" "app_sg" {
  name        = "coffeeshop-sg"
  description = "Security group for Coffeeshop services"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 9090
    to_port     = 9090
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}