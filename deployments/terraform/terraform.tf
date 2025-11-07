terraform {
  // The s3 bucket must be manually created outside of this 
  // terraform configuration. This is where terraform will
  // store the .tfstate file. It is recommended that you enable
  // versioning on the s3 bucket to be able to rollback your
  // terraform .tfstate file.
  backend "s3" {
    // You can use any name for the bucket that you like.
    bucket       = "zanesworld-terraform-state-files"
    key          = "gloomhaven-companion-service/terraform.tfstate"
    region       = "us-east-1"
    use_lockfile = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.17.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.7.1"
    }
    null = {
      source = "hashicorp/null"
    }
  }

  required_version = "~> 1.2"
}
