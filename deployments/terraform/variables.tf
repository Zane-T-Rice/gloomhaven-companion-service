# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

# Input variable definitions

variable "aws_region" {
  description = "AWS region for all resources."

  type    = string
  default = "us-east-1"
}

locals {
  archive_path = "../../dist/gloomhaven-companion-service.zip"
  binary_name  = "bootstrap"
  binary_path  = "../../dist/bootstrap"
  src_path     = "../../cmd/gloomhaven-companion-service"
}