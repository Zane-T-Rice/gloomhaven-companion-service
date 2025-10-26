variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "us-east-1"
}

locals {
  archive_path = "../../dist/gloomhaven-companion-service.zip"
  binary_name  = "bootstrap"
  binary_path  = "../../dist/bootstrap"
}