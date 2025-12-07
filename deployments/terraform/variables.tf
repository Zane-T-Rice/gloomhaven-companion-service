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
  gloomhaven_companion_service_lambda_execution_role = "AWSLambdaBasicExecutionRole"
  gloomhaven_companion_service_websocket_domain_name = "ws.zanesworld.click"
  gloomhaven_companion_service_api_domain_name       = "api.zanesworld.click"

  gloomhaven_companion_service_archive_path = "../../dist/gloomhaven-companion-service/gloomhaven-companion-service.zip"
  gloomhaven_companion_service_binary_name  = "bootstrap"
  gloomhaven_companion_service_binary_path  = "../../dist/gloomhaven-companion-service/bootstrap"

  gloomhaven_companion_service_websocket_connect_archive_path = "../../dist/gloomhaven-companion-service-websocket-connect/gloomhaven-companion-service-websocket-connect.zip"
  gloomhaven_companion_service_websocket_connect_binary_name  = "bootstrap"
  gloomhaven_companion_service_websocket_connect_binary_path  = "../../dist/gloomhaven-companion-service-websocket-connect/bootstrap"

  gloomhaven_companion_service_websocket_default_archive_path = "../../dist/gloomhaven-companion-service-websocket-default/gloomhaven-companion-service-websocket-default.zip"
  gloomhaven_companion_service_websocket_default_binary_name  = "bootstrap"
  gloomhaven_companion_service_websocket_default_binary_path  = "../../dist/gloomhaven-companion-service-websocket-default/bootstrap"

  gloomhaven_companion_service_websocket_disconnect_archive_path = "../../dist/gloomhaven-companion-service-websocket-disconnect/gloomhaven-companion-service-websocket-disconnect.zip"
  gloomhaven_companion_service_websocket_disconnect_binary_name  = "bootstrap"
  gloomhaven_companion_service_websocket_disconnect_binary_path  = "../../dist/gloomhaven-companion-service-websocket-disconnect/bootstrap"

  environment_variables = {
    // Auth0 Audience (Auth0 API Identifier)
    AUDIENCE = "https://zanesworld.click"

    // Auth0 tenant/issuer. Something like https://[-a-z1-9]+.us.auth0.com.
    ISSUER = "https://dev-7n448ak2gn3oqctx.us.auth0.com/"

    // This is the url for the service proper.
    // Either the custom domain+base path that you set up or
    // the invoke url from the API gateway -> Stages -> prod.
    GLOOMHAVEN_COMPANION_SERVICE_URL = "https://api.zanesworld.click/gloomhaven-companion-service"

    // The domain of the frontend to allow the origin for CORS.
    WEBSITE_DOMAIN = "https://zanesworld.click"

    // If you set up a custom domain, this would be the base path
    // for that domain.  By default, custom_domain.tf generates a
    // base path using the the service name /gloomhaven-companion-service.
    // Without a custom domain, this secret value is blank.
    API_GATEWAY_BASE_PATH = "/gloomhaven-companion-service"

    // This is the url for the websocket.
    // Either the custom domain+base path that you set up or
    // the WebSocket URL from the API gateway -> Stages -> prod.
    WEB_SOCKETS_URL = "https://ws.zanesworld.click/gloomhaven-companion-service"
  }
}