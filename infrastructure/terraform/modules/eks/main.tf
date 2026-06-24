# ── EKS Cluster ────────────────────────────────────────────────────

resource "aws_eks_cluster" "main" {
  name     = "sravas-${var.environment}"
  version  = var.cluster_version
  role_arn = var.cluster_role_arn

  vpc_config {
    subnet_ids              = var.subnet_ids
    security_group_ids      = var.cluster_security_group_ids
    endpoint_private_access = true
    endpoint_public_access  = var.endpoint_public_access
    public_access_cidrs     = var.endpoint_public_access ? var.public_access_cidrs : null
  }

  tags = merge({
    Name        = "sravas-${var.environment}"
    Environment = var.environment
  }, var.tags)
}

resource "aws_eks_addon" "vpc_cni" {
  cluster_name                = aws_eks_cluster.main.name
  addon_name                  = "vpc-cni"
  resolve_conflicts_on_create = "OVERWRITE"
  resolve_conflicts_on_update = "OVERWRITE"

  tags = merge({
    Name        = "sravas-vpc-cni-${var.environment}"
    Environment = var.environment
  }, var.tags)
}

resource "aws_eks_addon" "coredns" {
  cluster_name                = aws_eks_cluster.main.name
  addon_name                  = "coredns"
  resolve_conflicts_on_create = "OVERWRITE"
  resolve_conflicts_on_update = "OVERWRITE"

  tags = merge({
    Name        = "sravas-coredns-${var.environment}"
    Environment = var.environment
  }, var.tags)

  depends_on = [aws_eks_addon.vpc_cni]
}

resource "aws_eks_addon" "kube_proxy" {
  cluster_name                = aws_eks_cluster.main.name
  addon_name                  = "kube-proxy"
  resolve_conflicts_on_create = "OVERWRITE"
  resolve_conflicts_on_update = "OVERWRITE"

  tags = merge({
    Name        = "sravas-kube-proxy-${var.environment}"
    Environment = var.environment
  }, var.tags)
}

# ── Managed Node Group ─────────────────────────────────────────────

resource "aws_eks_node_group" "main" {
  cluster_name    = aws_eks_cluster.main.name
  node_group_name = "sravas-${var.environment}-nodes"
  node_role_arn   = var.node_role_arn
  subnet_ids      = var.subnet_ids
  instance_types  = var.instance_types

  scaling_config {
    desired_size = var.desired_size
    min_size     = var.min_size
    max_size     = var.max_size
  }

  update_config {
    max_unavailable = 1
  }

  tags = merge({
    Name        = "sravas-node-group-${var.environment}"
    Environment = var.environment
  }, var.tags)

  depends_on = [
    aws_eks_addon.vpc_cni,
    aws_eks_addon.coredns,
    aws_eks_addon.kube_proxy,
  ]
}

# ── OIDC Provider for IRSA ────────────────────────────────────────

data "tls_certificate" "eks" {
  url = aws_eks_cluster.main.identity[0].oidc[0].issuer
}

resource "aws_iam_openid_connect_provider" "eks" {
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = [data.tls_certificate.eks.certificates[0].sha1_fingerprint]
  url             = aws_eks_cluster.main.identity[0].oidc[0].issuer
}

# ── Cluster Autoscaler IRSA ────────────────────────────────────────

locals {
  oidc_provider    = replace(aws_eks_cluster.main.identity[0].oidc[0].issuer, "https://", "")
  autoscaler_name  = "sravas-cluster-autoscaler-${var.environment}"
}

resource "aws_iam_role" "cluster_autoscaler" {
  name = local.autoscaler_name

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = aws_iam_openid_connect_provider.eks.arn
      }
      Action = "sts:AssumeRoleWithWebIdentity"
      Condition = {
        StringEquals = {
          "${local.oidc_provider}:sub" = "system:serviceaccount:kube-system:cluster-autoscaler"
          "${local.oidc_provider}:aud" = "sts.amazonaws.com"
        }
      }
    }]
  })

  tags = merge({
    Name        = local.autoscaler_name
    Environment = var.environment
  }, var.tags)
}

resource "aws_iam_policy" "cluster_autoscaler" {
  name        = local.autoscaler_name
  description = "Cluster autoscaler permissions for EKS"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "autoscaling:DescribeAutoScalingGroups",
        "autoscaling:DescribeAutoScalingInstances",
        "autoscaling:DescribeLaunchConfigurations",
        "autoscaling:DescribeScalingActivities",
        "autoscaling:SetDesiredCapacity",
        "autoscaling:TerminateInstanceInAutoScalingGroup",
        "ec2:DescribeLaunchTemplateVersions",
        "ec2:DescribeInstanceTypes",
        "ec2:DescribeImages",
      ]
      Resource = ["*"]
    }]
  })
}

resource "aws_iam_role_policy_attachment" "cluster_autoscaler" {
  policy_arn = aws_iam_policy.cluster_autoscaler.arn
  role       = aws_iam_role.cluster_autoscaler.name
}
