// NOTE:
//
// There are a few resources that need to be created
// or updated manually.
// 
// There is more information about these in
// custom-domains.tf, data.tf, secrets.tf, and terraform.tf.

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      zanesworld = "gloomhaven-companion-service"
    }
  }
}

