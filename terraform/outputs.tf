output "public_ip" {
  description = "Public IP address of the EC2 instance"
  value       = aws_instance.app_server.public_ip
}

output "ssh_connection_command" {
  description = "Command to connect to the instance via SSH"
  value       = "ssh -i '${var.key_name}.pem' ubuntu@${aws_instance.app_server.public_ip}"
}

output "grafana_url" {
  description = "URL for Grafana dashboard"
  value       = "http://${aws_instance.app_server.public_ip}:3000"
}

output "prometheus_url" {
  description = "URL for Prometheus dashboard"
  value       = "http://${aws_instance.app_server.public_ip}:9090"
}