// All of these secrets will need to be set manually in AWS Secrets Manager
// after they are created by Terraform.

// If you set up a custom domain, this would be the base path
// for that domain.  By default, custom_domain.tf generates a 
// base path using the the service name /gloomhaven-companion-service.
// Without a custom domain, this secret value is blank.
resource "aws_secretsmanager_secret" "api_gateway_base_path" {
  name = "gloomhaven-companion-service-api-gateway-base-path"
}
// Auth0 Audience (Auth0 API Identifier)
resource "aws_secretsmanager_secret" "audience" {
  name = "gloomhaven-companion-service-audience"
}
// Auth0 tenant/issuer. Something like https://[-a-z1-9]+.us.auth0.com.
resource "aws_secretsmanager_secret" "issuer" {
  name = "gloomhaven-companion-service-issuer"
}
// This is the url for the service proper.
// Either the custom domain+base path that you set up or
// the invoke url from the API gateway -> Stages -> prod.
resource "aws_secretsmanager_secret" "url" {
  name = "gloomhaven-companion-service-url"
}
// The domain of the frontend to allow the origin for CORS.
resource "aws_secretsmanager_secret" "website_domain" {
  name = "gloomhaven-companion-service-website-domain"
}
// This is the url for the websocket.
// Either the custom domain+base path that you set up or
// the WebSocket URL from the API gateway -> Stages -> prod.
resource "aws_secretsmanager_secret" "web_sockets_url" {
  name = "gloomhaven-companion-service-websockets-url"
}