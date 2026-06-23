provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "sravas"
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}
