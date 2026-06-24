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
