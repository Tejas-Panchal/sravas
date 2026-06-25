environment                 = "prod"
availability_zones          = ["us-east-1a", "us-east-1b", "us-east-1c"]
single_nat_gateway          = false

rds_multi_az               = true
rds_deletion_protection     = true
rds_skip_final_snapshot     = false
rds_backup_retention_period = 30
