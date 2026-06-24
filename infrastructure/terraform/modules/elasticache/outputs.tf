output "cache_endpoint" {
  description = "Primary endpoint address for the Redis replication group"
  value       = aws_elasticache_replication_group.main.primary_endpoint_address
}

output "cache_reader_endpoint" {
  description = "Reader endpoint address for the Redis replication group"
  value       = aws_elasticache_replication_group.main.reader_endpoint_address
}

output "cache_port" {
  description = "Redis port number"
  value       = aws_elasticache_replication_group.main.port
}

output "cache_endpoint_port" {
  description = "Redis endpoint:port string"
  value       = format("%s:%d", aws_elasticache_replication_group.main.primary_endpoint_address, aws_elasticache_replication_group.main.port)
}

output "cache_replication_group_id" {
  description = "Replication group ID"
  value       = aws_elasticache_replication_group.main.id
}

output "cache_subnet_group_name" {
  description = "Name of the ElastiCache subnet group"
  value       = aws_elasticache_subnet_group.main.name
}
