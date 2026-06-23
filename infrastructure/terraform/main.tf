provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "sravas"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }

  assume_role {
    role_arn     = var.assume_role_arn
    session_name = "sravas-tf-${var.environment}"
  }
}
