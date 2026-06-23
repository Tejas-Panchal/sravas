terraform {
  backend "s3" {
    bucket         = "sravas-tfstate-a86b1c95"
    key            = "sravas/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    use_lockfile   = true
  }
}
