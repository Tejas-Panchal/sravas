module "rds" {
  source = "./modules/rds"

  environment              = var.environment
  subnet_ids               = module.vpc.private_subnet_ids
  vpc_id                   = module.vpc.vpc_id
  security_group_ids       = [aws_security_group.db.id]
  multi_az                 = var.rds_multi_az
  deletion_protection      = var.rds_deletion_protection
  skip_final_snapshot      = var.rds_skip_final_snapshot
  backup_retention_period  = var.rds_backup_retention_period
}
