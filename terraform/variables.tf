variable "instance_type" { default = "c7i-flex.large" }
variable "ami_id"        { default = "ami-0014ce3e52359afbd" } 
variable "key_name"      { default = "coffeeshop-aws" }

variable "aws_access_key" {
  description = "AWS Access Key"
  type        = string
}

variable "aws_secret_key" {
  description = "AWS Secret Key"
  type        = string
}