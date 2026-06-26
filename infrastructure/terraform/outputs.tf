output "tfstate_bucket" {
  description = "S3 bucket name containing Terraform state"
  value       = aws_s3_bucket.tfstate.bucket
}

output "tfstate_lock_table" {
  description = "DynamoDB table name for Terraform state locking"
  value       = aws_dynamodb_table.lock.name
}

output "kubeconfig_command" {
  description = "AWS CLI command to configure kubectl for the EKS cluster"
  value       = format("aws eks update-kubeconfig --region %s --name %s", var.aws_region, module.eks.cluster_id)
}

output "cloudfront_domain_name" {
  description = "Domain name of the CloudFront distribution (use for CNAME / app URL)"
  value       = module.cloudfront.cloudfront_domain_name
}

output "cloudfront_distribution_id" {
  description = "ID of the CloudFront distribution (for cache invalidation)"
  value       = module.cloudfront.cloudfront_distribution_id
}

output "alb_controller_role_arn" {
  description = "ARN of the IAM role for the AWS Load Balancer Controller (use in IRSA annotation)"
  value       = module.eks.alb_controller_role_arn
}

output "db_endpoint" {
  description = "RDS endpoint (host:port) — use in service ConfigMaps"
  value       = module.rds.db_endpoint
}

output "redis_endpoint" {
  description = "Redis primary endpoint — use in service ConfigMaps"
  value       = module.elasticache.cache_endpoint
}

output "redis_port" {
  description = "Redis port number"
  value       = module.elasticache.cache_port
}

output "upload_bucket_id" {
  description = "Upload S3 bucket name — use in upload-service ConfigMap"
  value       = module.s3.upload_bucket_id
}

output "db_username" {
  description = "RDS master username — use in service Secrets"
  value       = module.rds.db_username
}

output "db_password" {
  description = "RDS master password — use in service Secrets"
  value       = module.rds.db_password
  sensitive   = true
}
