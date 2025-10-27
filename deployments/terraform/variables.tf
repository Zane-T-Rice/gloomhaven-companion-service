variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "us-east-1"
}

variable "dynamodb_table_name" {
  description = "DynamoDB table name for the gloomhaven companion service."

  type    = string
  default = "gloomhaven-companion-service"
}

locals {
  archive_path = "../../dist/gloomhaven-companion-service.zip"
  binary_name  = "bootstrap"
  binary_path  = "../../dist/bootstrap"
}