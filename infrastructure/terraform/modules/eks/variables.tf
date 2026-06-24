variable "environment" {
  description = "Deployment environment (dev / prod)"
  type        = string
}

variable "cluster_role_arn" {
  description = "IAM role ARN for the EKS control plane"
  type        = string
}

variable "node_role_arn" {
  description = "IAM role ARN for the EKS managed node group"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet IDs to deploy the EKS cluster into"
  type        = list(string)
}

variable "cluster_security_group_ids" {
  description = "Additional security group IDs to attach to the EKS cluster"
  type        = list(string)
  default     = []
}

variable "cluster_version" {
  description = "Kubernetes version for the EKS cluster"
  type        = string
  default     = "1.30"
}

variable "instance_types" {
  description = "EC2 instance types for the managed node group"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "desired_size" {
  description = "Desired number of nodes in the managed node group"
  type        = number
  default     = 2
}

variable "min_size" {
  description = "Minimum number of nodes in the managed node group"
  type        = number
  default     = 2
}

variable "max_size" {
  description = "Maximum number of nodes in the managed node group"
  type        = number
  default     = 4
}

variable "endpoint_public_access" {
  description = "Enable public access to the EKS API server"
  type        = bool
  default     = true
}

variable "public_access_cidrs" {
  description = "CIDR blocks allowed to access the EKS API server publicly"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "tags" {
  description = "Additional tags for all EKS resources"
  type        = map(string)
  default     = {}
}
