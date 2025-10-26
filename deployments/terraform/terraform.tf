terraform {

  # These values have to be hardcoded and Terraform cannot generate this S3 bucket
  # as well as use it as a state storage device. The S3 bucket has versioning
  # turned on to allow rollbacks of Terraform state.
  backend "s3" {
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
