resource "aws_secretsmanager_secret" "api_gateway_base_path" {
  name = "gloomhaven-companion-service-api-gateway-base-path"
}
resource "aws_secretsmanager_secret" "audience" {
  name = "gloomhaven-companion-service-audience"
}
resource "aws_secretsmanager_secret" "issuer" {
  name = "gloomhaven-companion-service-issuer"
}
resource "aws_secretsmanager_secret" "url" {
  name = "gloomhaven-companion-service-url"
}
resource "aws_secretsmanager_secret" "website_domain" {
  name = "gloomhaven-companion-service-website-domain"
}
resource "aws_secretsmanager_secret" "web_sockets_url" {
  name = "gloomhaven-companion-service-websockets-url"
}