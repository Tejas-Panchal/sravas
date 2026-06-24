variable "environment" {
  description = "Deployment environment (dev / prod)"
  type        = string
}

variable "cloudfront_oai_iam_arn" {
  description = "IAM ARN of the CloudFront Origin Access Identity"
  type        = string
}

variable "cors_allowed_origins" {
  description = "Origins allowed to make CORS requests to the upload bucket"
  type        = list(string)
  default     = ["*"]
}

variable "transition_ia_days" {
  description = "Days before transitioning objects to STANDARD_IA"
  type        = number
  default     = 30
}

variable "transition_glacier_days" {
  description = "Days before transitioning objects to GLACIER_IR"
  type        = number
  default     = 90
}

variable "expiration_days" {
  description = "Days before objects are permanently deleted"
  type        = number
  default     = 365
}

variable "log_expiration_days" {
  description = "Days before access log objects are permanently deleted"
  type        = number
  default     = 90
}

variable "tags" {
  description = "Additional tags for all S3 resources"
  type        = map(string)
  default     = {}
}
