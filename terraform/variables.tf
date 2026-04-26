variable "aws_access_key"{
    description ="AWS Access key"
    type = string
}

variable "aws_secret_key" {
  description = "AWS Secret Key"
  type        = string
  sensitive   = true 
}