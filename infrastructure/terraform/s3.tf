module "s3" {
  source = "./modules/s3"

  environment            = var.environment
  cloudfront_oai_iam_arn = module.iam.cloudfront_oai_iam_arn
}
