variable "key_name" {
  description = "Name of the existing EC2 Key Pair for SSH access"
  type        = string
  default     = "coffeeshop-aws"
}

variable "instance_type" {
  description = "EC2 Instance Type"
  type        = string
  default     = "c7i-flex.large"
}

variable "ami_id" {
  description = "Ubuntu AMI ID for eu-north-1"
  type        = string
  default     = "ami-0014ce3e52359afbd"
} 

variable "aws_access_key" {
  description = "AWS Access Key"
  type        = string
}

variable "aws_secret_key" {
  description = "AWS Secret Key"
  type        = string
}

variable "github_token" {
  description = "GitHub Personal Access Token"
  type        = string
  sensitive   = true
  default     = ""
}