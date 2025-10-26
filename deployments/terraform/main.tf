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

// zip the binary, as we can use only zip files to AWS lambda
data "archive_file" "function_archive" {
  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

// Secrets for the service
resource "aws_secretsmanager_secret" "audience" {
  name = "gloomhaven-companion-service-audience"
}
resource "aws_secretsmanager_secret" "issuer" {
  name = "gloomhaven-companion-service-issuer"
}

// create the lambda function from zip file
resource "aws_lambda_function" "gloomhaven-companion-service" {
  function_name = "gloomhaven-companion-service"
  description   = "Backend service for the gloomhaven companion."
  role          = data.aws_iam_role.lambda_exec.arn
  handler       = local.binary_name
  memory_size   = 128

  filename         = local.archive_path
  source_code_hash = data.archive_file.function_archive.output_base64sha256

  runtime = "provided.al2023"
}
resource "aws_lambda_function_url" "gloomhaven-companion-service" {
  function_name      = aws_lambda_function.gloomhaven-companion-service.id
  authorization_type = "NONE"
  cors {
    allow_origins = ["*"]
  }
}

// This is the target permission that needs to be added to the function's policy 
// However, there is no way to add the Condition in terraform.
// For now the policy is added manually to the function.
// {
//   "Sid": "FunctionURLAllowInvokeAction",
//   "Effect": "Allow",
//   "Principal": "*",
//   "Action": "lambda:InvokeFunction",
//   "Resource": "arn:aws:lambda:us-east-1:675101127982:function:gloomhaven-companion-service",
//   "Condition": {
//     "Bool": {
//       "lambda:InvokedViaFunctionUrl": "true"
//     }
//   }
// }
// resource "aws_lambda_permission" "allow_public_url_invocation" {
//   statement_id = "FunctionURLAllowInvokeAction"
//   action       = "lambda:InvokeFunction"
//   function_name = aws_lambda_function.gloomhaven-companion-service.function_name
//   principal     = "*"
// 
//   depends_on = [aws_lambda_function_url.gloomhaven-companion-service]
// }


resource "aws_cloudwatch_log_group" "gloomhaven-companion-service" {
  name = "/aws/lambda/${aws_lambda_function.gloomhaven-companion-service.function_name}"

  retention_in_days = 30
}



