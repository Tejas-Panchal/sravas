output "upload_bucket_id" {
  description = "ID of the upload S3 bucket"
  value       = aws_s3_bucket.upload.id
}

output "upload_bucket_arn" {
  description = "ARN of the upload S3 bucket"
  value       = aws_s3_bucket.upload.arn
}

output "upload_bucket_regional_domain_name" {
  description = "Regional domain name of the upload bucket (for CloudFront origin)"
  value       = aws_s3_bucket.upload.bucket_regional_domain_name
}

output "logs_bucket_id" {
  description = "ID of the access logs S3 bucket"
  value       = aws_s3_bucket.logs.id
}

output "logs_bucket_arn" {
  description = "ARN of the access logs S3 bucket"
  value       = aws_s3_bucket.logs.arn
}
