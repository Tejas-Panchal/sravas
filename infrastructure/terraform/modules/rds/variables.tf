variable "environment" {
  description = "Deployment environment (dev / prod)"
  type        = string
}

variable "subnet_ids" {
  description = "Private subnet IDs for the DB subnet group"
  type        = list(string)
}

variable "vpc_id" {
  description = "VPC ID for the parameter group"
  type        = string
}

variable "security_group_ids" {
  description = "Security group IDs to attach to the RDS instance"
  type        = list(string)
}

variable "db_name" {
  description = "Name of the PostgreSQL database to create"
  type        = string
  default     = "sravas"
}

variable "db_username" {
  description = "Master username for the RDS instance"
  type        = string
  default     = "sravas"
}

variable "db_password" {
  description = "Master password for RDS (empty = auto-generated)"
  type        = string
  default     = ""
  sensitive   = true
}

variable "instance_class" {
  description = "RDS instance type"
  type        = string
  default     = "db.t3.small"
}

variable "multi_az" {
  description = "Enable multi-AZ deployment (set true for prod)"
  type        = bool
  default     = false
}

variable "allocated_storage" {
  description = "Allocated storage in GB"
  type        = number
  default     = 20
}

variable "max_allocated_storage" {
  description = "Maximum storage for autoscaling (0 to disable)"
  type        = number
  default     = 100
}

variable "engine_version" {
  description = "PostgreSQL engine version"
  type        = string
  default     = "16.3"
}

variable "backup_retention_period" {
  description = "Number of days to retain automated backups"
  type        = number
  default     = 7
}

variable "backup_window" {
  description = "Preferred backup window (UTC)"
  type        = string
  default     = "03:00-04:00"
}

variable "maintenance_window" {
  description = "Preferred maintenance window (UTC)"
  type        = string
  default     = "sun:04:00-sun:05:00"
}

variable "deletion_protection" {
  description = "Enable deletion protection (set true for prod)"
  type        = bool
  default     = false
}

variable "skip_final_snapshot" {
  description = "Skip final snapshot on deletion (set false for prod)"
  type        = bool
  default     = true
}

variable "publicly_accessible" {
  description = "Allow public internet access to the instance"
  type        = bool
  default     = false
}

variable "storage_encrypted" {
  description = "Enable storage encryption at rest"
  type        = bool
  default     = true
}

variable "tags" {
  description = "Additional tags for all RDS resources"
  type        = map(string)
  default     = {}
}
