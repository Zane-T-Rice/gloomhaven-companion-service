// zip the binary, as we can use only zip files to AWS lambda
data "archive_file" "function_archive" {
  type        = "zip"
  source_file = local.gloomhaven_companion_service_binary_path
  output_path = local.gloomhaven_companion_service_archive_path
}

// create the lambda function from zip file
resource "aws_lambda_function" "gloomhaven_companion_service" {
  function_name = "gloomhaven-companion-service"
  description   = "Backend service for the gloomhaven companion."
  role          = data.aws_iam_role.lambda_exec.arn
  handler       = local.gloomhaven_companion_service_binary_name
  memory_size   = 128

  filename         = local.gloomhaven_companion_service_archive_path
  source_code_hash = data.archive_file.function_archive.output_base64sha256

  runtime = "provided.al2023"
}

resource "aws_cloudwatch_log_group" "gloomhaven_companion_service" {
  name = "/aws/lambda/${aws_lambda_function.gloomhaven_companion_service.function_name}"

  retention_in_days = 30
}
