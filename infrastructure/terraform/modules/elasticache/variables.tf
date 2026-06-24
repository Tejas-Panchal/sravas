variable "environment" {
  description = "Deployment environment (dev / prod)"
  type        = string
}

variable "subnet_ids" {
  description = "Private subnet IDs for the ElastiCache subnet group"
  type        = list(string)
}

variable "security_group_ids" {
  description = "Security group IDs to attach to the replication group"
  type        = list(string)
}

variable "node_type" {
  description = "ElastiCache node instance type"
  type        = string
  default     = "cache.t3.micro"
}

variable "engine_version" {
  description = "Redis engine version"
  type        = string
  default     = "7.1"
}

variable "num_cache_clusters" {
  description = "Number of cache clusters (used when cluster_mode_enabled = false)"
  type        = number
  default     = 1
}

variable "cluster_mode_enabled" {
  description = "Enable cluster mode (set true for prod)"
  type        = bool
  default     = false
}

variable "num_node_groups" {
  description = "Number of shards in cluster mode"
  type        = number
  default     = 2
}

variable "replicas_per_node_group" {
  description = "Number of replica nodes per shard in cluster mode"
  type        = number
  default     = 1
}

variable "port" {
  description = "Redis port number"
  type        = number
  default     = 6379
}

variable "automatic_failover" {
  description = "Enable automatic failover (required for multi-AZ / cluster mode)"
  type        = bool
  default     = false
}

variable "multi_az_enabled" {
  description = "Enable multi-AZ replication"
  type        = bool
  default     = false
}

variable "tags" {
  description = "Additional tags for all ElastiCache resources"
  type        = map(string)
  default     = {}
}
