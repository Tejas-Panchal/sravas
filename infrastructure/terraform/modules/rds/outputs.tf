output "db_hostname" {
  description = "RDS instance hostname (without port)"
  value       = aws_db_instance.main.address
}

output "db_port" {
  description = "RDS instance port"
  value       = aws_db_instance.main.port
}

output "db_endpoint" {
  description = "RDS instance endpoint (hostname:port)"
  value       = format("%s:%s", aws_db_instance.main.address, aws_db_instance.main.port)
}

output "db_name" {
  description = "Name of the database"
  value       = aws_db_instance.main.db_name
}

output "db_username" {
  description = "Master username"
  value       = aws_db_instance.main.username
}

output "db_password" {
  description = "Master password (auto-generated if not provided)"
  value       = local.db_password
  sensitive   = true
}

output "db_arn" {
  description = "ARN of the RDS instance"
  value       = aws_db_instance.main.arn
}

output "db_instance_id" {
  description = "ID of the RDS instance"
  value       = aws_db_instance.main.id
}

output "db_subnet_group_name" {
  description = "Name of the DB subnet group"
  value       = aws_db_subnet_group.main.name
}
