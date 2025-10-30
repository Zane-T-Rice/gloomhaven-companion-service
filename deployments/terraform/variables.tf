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
  gloomhaven_companion_service_archive_path           = "../../dist/gloomhaven-companion-service/gloomhaven-companion-service.zip"
  gloomhaven_companion_service_binary_name            = "bootstrap"
  gloomhaven_companion_service_binary_path            = "../../dist/gloomhaven-companion-service/bootstrap"
  gloomahven_companion_service_websocket_archive_path = "../../dist/gloomhaven-companion-service-websocket/gloomhaven-companion-service-websocket.zip"
  gloomahven_companion_service_websocket_binary_name  = "bootstrap"
  gloomahven_companion_service_websocket_binary_path  = "../../dist/gloomhaven-companion-service-websocket/bootstrap"
}