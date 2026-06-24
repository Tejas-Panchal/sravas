data "aws_caller_identity" "current" {}

# ── EKS Cluster Role ──────────────────────────────────────────────

resource "aws_iam_role" "eks_cluster" {
  name = "sravas-eks-cluster-role-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = { Service = "eks.amazonaws.com" }
      Action = "sts:AssumeRole"
    }]
  })

  tags = merge({
    Name        = "sravas-eks-cluster-role-${var.environment}"
    Environment = var.environment
  }, var.tags)
}

resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_cluster.name
}

resource "aws_iam_role_policy_attachment" "eks_vpc_resource_controller" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
  role       = aws_iam_role.eks_cluster.name
}

# ── Node Instance Role ────────────────────────────────────────────

resource "aws_iam_role" "eks_node" {
  name = "sravas-eks-node-role-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
      Action = "sts:AssumeRole"
    }]
  })

  tags = merge({
    Name        = "sravas-eks-node-role-${var.environment}"
    Environment = var.environment
  }, var.tags)
}

resource "aws_iam_role_policy_attachment" "eks_worker_node" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "eks_cni" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_node.name
}

resource "aws_iam_role_policy_attachment" "ecr_readonly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.eks_node.name
}

# ── S3 Upload Access Policy ───────────────────────────────────────

resource "aws_iam_policy" "s3_upload_access" {
  name        = "sravas-s3-upload-access-${var.environment}"
  description = "Allow read/write to the upload S3 bucket"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "S3UploadReadWrite"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
        ]
        Resource = [
          "arn:aws:s3:::sravas-uploads-${var.environment}",
          "arn:aws:s3:::sravas-uploads-${var.environment}/*"
        ]
      },
      {
        Sid    = "S3UploadListBucket"
        Effect = "Allow"
        Action = ["s3:ListBucket"]
        Resource = [
          "arn:aws:s3:::sravas-uploads-${var.environment}"
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "s3_upload_access" {
  policy_arn = aws_iam_policy.s3_upload_access.arn
  role       = aws_iam_role.eks_node.name
}

# ── CloudFront Origin Access Identity ─────────────────────────────

resource "aws_cloudfront_origin_access_identity" "main" {
  comment = "sravas-${var.environment}"
}
