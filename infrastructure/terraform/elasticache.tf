module "elasticache" {
  source = "./modules/elasticache"

  environment        = var.environment
  subnet_ids         = module.vpc.private_subnet_ids
  security_group_ids = [aws_security_group.cache.id]
}
