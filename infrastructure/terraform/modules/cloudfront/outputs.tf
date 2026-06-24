output "cloudfront_domain_name" {
  description = "Domain name of the CloudFront distribution"
  value       = aws_cloudfront_distribution.main.domain_name
}

output "cloudfront_distribution_id" {
  description = "ID of the CloudFront distribution (for cache invalidation)"
  value       = aws_cloudfront_distribution.main.id
}

output "cloudfront_hosted_zone_id" {
  description = "Hosted zone ID for Route53 alias records"
  value       = aws_cloudfront_distribution.main.hosted_zone_id
}
