output "eks_cluster_role_arn" {
  description = "ARN of the EKS cluster IAM role"
  value       = aws_iam_role.eks_cluster.arn
}

output "eks_node_role_arn" {
  description = "ARN of the EKS node instance IAM role"
  value       = aws_iam_role.eks_node.arn
}

output "s3_upload_policy_arn" {
  description = "ARN of the S3 upload access policy"
  value       = aws_iam_policy.s3_upload_access.arn
}

output "cloudfront_oai_id" {
  description = "ID of the CloudFront Origin Access Identity"
  value       = aws_cloudfront_origin_access_identity.main.id
}

output "cloudfront_oai_iam_arn" {
  description = "IAM ARN of the CloudFront Origin Access Identity (use in S3 bucket policy)"
  value       = aws_cloudfront_origin_access_identity.main.iam_arn
}
