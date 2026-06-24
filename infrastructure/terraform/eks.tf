module "eks" {
  source = "./modules/eks"

  environment                = var.environment
  cluster_role_arn           = module.iam.eks_cluster_role_arn
  node_role_arn              = module.iam.eks_node_role_arn
  subnet_ids                 = module.vpc.private_subnet_ids
  cluster_security_group_ids = [aws_security_group.internal.id]
}
