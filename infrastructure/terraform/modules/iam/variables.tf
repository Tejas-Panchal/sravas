variable "environment" {
  description = "Deployment environment (dev / prod)"
  type        = string
}

variable "tags" {
  description = "Additional tags for all IAM resources"
  type        = map(string)
  default     = {}
}
