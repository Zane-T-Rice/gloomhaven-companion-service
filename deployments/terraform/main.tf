provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      zanesworld = "gloomhaven-companion-service"
    }
  }
}

# Use existing Role to avoid needing to give Software Engineer Group
# Full IAM Access.
data "aws_iam_role" "lambda_exec" {
  name = "AWSLambdaBasicExecutionRole"
}

// Secrets for the service
resource "aws_secretsmanager_secret" "audience" {
  name = "gloomhaven-companion-service-audience"
}
resource "aws_secretsmanager_secret" "issuer" {
  name = "gloomhaven-companion-service-issuer"
}
resource "aws_secretsmanager_secret" "url" {
  name = "gloomhaven-companion-service-url"
}
