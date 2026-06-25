variable "aws_region" {
  description = "AWS region to deploy into"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Deployment environment (dev / prod)"
  type        = string
  default     = "dev"
}

variable "assume_role_arn" {
  description = "IAM role ARN to assume for Terraform (empty = use current credentials)"
  type        = string
  default     = ""
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "AWS availability zones to deploy subnets into"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b", "us-east-1c"]
}

variable "single_nat_gateway" {
  description = "Use a single NAT Gateway for all subnets (cheaper, suitable for dev)"
  type        = bool
  default     = true
}

variable "rds_multi_az" {
  description = "Enable multi-AZ deployment for RDS (set true for prod)"
  type        = bool
  default     = false
}

variable "rds_deletion_protection" {
  description = "Enable deletion protection for RDS (set true for prod)"
  type        = bool
  default     = false
}

variable "rds_skip_final_snapshot" {
  description = "Skip final snapshot on RDS deletion (set false for prod)"
  type        = bool
  default     = true
}

variable "rds_backup_retention_period" {
  description = "Number of days to retain RDS automated backups"
  type        = number
  default     = 7
}

variable "elasticache_cluster_mode_enabled" {
  description = "Enable ElastiCache cluster mode (set true for prod)"
  type        = bool
  default     = false
}
