variable "environment" {
  description = "Deployment environment (dev / prod)"
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "AWS availability zones to deploy subnets into"
  type        = list(string)
}

variable "single_nat_gateway" {
  description = "Use a single NAT Gateway for all subnets (cheaper)"
  type        = bool
  default     = true
}

variable "enable_dns_hostnames" {
  description = "Enable DNS hostnames in the VPC"
  type        = bool
  default     = true
}

variable "enable_dns_support" {
  description = "Enable DNS support in the VPC"
  type        = bool
  default     = true
}

variable "tags" {
  description = "Additional tags for all VPC resources"
  type        = map(string)
  default     = {}
}
