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
