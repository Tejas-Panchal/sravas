# ── Subnet Group ─────────────────────────────────────────────────

resource "aws_elasticache_subnet_group" "main" {
  name        = "sravas-${var.environment}"
  description = "Private subnets for ElastiCache ${var.environment}"
  subnet_ids  = var.subnet_ids

  tags = merge({
    Name        = "sravas-elasticache-subnet-${var.environment}"
    Environment = var.environment
  }, var.tags)
}

# ── Replication Group (Redis 7) ─────────────────────────────────

resource "aws_elasticache_replication_group" "main" {
  replication_group_id = "sravas-${var.environment}"
  description          = "Redis ${var.environment}"

  engine               = "redis"
  engine_version       = var.engine_version
  node_type            = var.node_type
  port                 = var.port
  cluster_mode         = var.cluster_mode_enabled ? "enabled" : "disabled"
  num_cache_clusters   = var.cluster_mode_enabled ? null : var.num_cache_clusters

  subnet_group_name          = aws_elasticache_subnet_group.main.name
  security_group_ids         = var.security_group_ids
  automatic_failover_enabled = var.automatic_failover
  multi_az_enabled           = var.multi_az_enabled

  snapshot_retention_limit = 1
  snapshot_window          = "03:00-04:00"

  tags = merge({
    Name        = "sravas-${var.environment}"
    Environment = var.environment
  }, var.tags)
}
