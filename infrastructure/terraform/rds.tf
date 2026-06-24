module "rds" {
  source = "./modules/rds"

  environment        = var.environment
  subnet_ids         = module.vpc.private_subnet_ids
  vpc_id             = module.vpc.vpc_id
  security_group_ids = [aws_security_group.db.id]
}
