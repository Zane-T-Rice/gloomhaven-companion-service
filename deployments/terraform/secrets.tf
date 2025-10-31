resource "aws_secretsmanager_secret" "audience" {
  name = "gloomhaven-companion-service-audience"
}
resource "aws_secretsmanager_secret" "issuer" {
  name = "gloomhaven-companion-service-issuer"
}
resource "aws_secretsmanager_secret" "url" {
  name = "gloomhaven-companion-service-url"
}
resource "aws_secretsmanager_secret" "website-domain" {
  name = "gloomhaven-companion-service-website-domain"
}