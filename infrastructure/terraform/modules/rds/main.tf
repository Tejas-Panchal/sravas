# ── Random Password (conditional) ─────────────────────────────────

resource "random_password" "master" {
  count   = var.db_password == "" ? 1 : 0
  length  = 24
  special = false
}

locals {
  db_password = var.db_password != "" ? var.db_password : random_password.master[0].result
}

# ── Parameter Group ──────────────────────────────────────────────

resource "aws_db_parameter_group" "main" {
  name   = "sravas-postgres16-${var.environment}"
  family = "postgres16"

  parameter {
    name  = "log_min_duration_statement"
    value = "1000"
  }

  parameter {
    name         = "rds.force_ssl"
    value        = "0"
    apply_method = "pending-reboot"
  }

  tags = merge({
    Name        = "sravas-postgres16-${var.environment}"
    Environment = var.environment
  }, var.tags)
}

# ── Subnet Group ─────────────────────────────────────────────────

resource "aws_db_subnet_group" "main" {
  name        = "sravas-${var.environment}"
  description = "Private subnets for RDS ${var.environment}"
  subnet_ids  = var.subnet_ids

  tags = merge({
    Name        = "sravas-db-subnet-${var.environment}"
    Environment = var.environment
  }, var.tags)
}

# ── RDS Instance ─────────────────────────────────────────────────

resource "aws_db_instance" "main" {
  identifier     = "sravas-${var.environment}"
  engine         = "postgres"
  engine_version = var.engine_version

  instance_class    = var.instance_class
  allocated_storage = var.allocated_storage
  storage_type      = "gp3"
  storage_encrypted = var.storage_encrypted

  db_name  = var.db_name
  username = var.db_username
  password = local.db_password

  multi_az             = var.multi_az
  publicly_accessible  = var.publicly_accessible
  parameter_group_name = aws_db_parameter_group.main.name
  db_subnet_group_name = aws_db_subnet_group.main.name
  vpc_security_group_ids = var.security_group_ids

  backup_retention_period = var.backup_retention_period
  backup_window           = var.backup_window
  maintenance_window      = var.maintenance_window

  auto_minor_version_upgrade = true
  copy_tags_to_snapshot      = true
  deletion_protection        = var.deletion_protection
  skip_final_snapshot       = var.skip_final_snapshot
  final_snapshot_identifier = var.skip_final_snapshot ? null : "sravas-${var.environment}-final-snapshot"

  tags = merge({
    Name        = "sravas-${var.environment}"
    Environment = var.environment
  }, var.tags)
}
