module "cloudfront" {
  source = "./modules/cloudfront"

  providers = {
    aws.us_east_1 = aws.us_east_1
  }

  environment                    = var.environment
  s3_bucket_regional_domain_name = module.s3.upload_bucket_regional_domain_name
  s3_bucket_id                   = module.s3.upload_bucket_id
  cloudfront_oai_path            = module.iam.cloudfront_oai_path
}
