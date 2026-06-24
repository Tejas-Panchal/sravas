variable "environment" {
  description = "Deployment environment (dev / prod)"
  type        = string
}

variable "s3_bucket_regional_domain_name" {
  description = "Regional domain name of the upload S3 bucket"
  type        = string
}

variable "s3_bucket_id" {
  description = "ID of the upload S3 bucket"
  type        = string
}

variable "cloudfront_oai_path" {
  description = "CloudFront origin access identity path for S3 origin config"
  type        = string
}

variable "alb_domain_name" {
  description = "ALB DNS name for the frontend origin (null = no ALB origin)"
  type        = string
  default     = null
}

variable "price_class" {
  description = "CloudFront price class (PriceClass_100 = US+EU, PriceClass_200 = +Asia, PriceClass_All = everywhere)"
  type        = string
  default     = "PriceClass_100"
}

variable "default_ttl" {
  description = "Default TTL for cached responses (seconds)"
  type        = number
  default     = 86400
}

variable "max_ttl" {
  description = "Maximum TTL for cached responses (seconds)"
  type        = number
  default     = 604800
}

variable "min_ttl" {
  description = "Minimum TTL for cached responses (seconds)"
  type        = number
  default     = 0
}

variable "enabled" {
  description = "Enable the CloudFront distribution"
  type        = bool
  default     = true
}

variable "tags" {
  description = "Additional tags for all CloudFront resources"
  type        = map(string)
  default     = {}
}
