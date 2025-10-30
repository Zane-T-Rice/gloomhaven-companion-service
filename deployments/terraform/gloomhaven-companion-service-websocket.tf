// zip the binary, as we can use only zip files to AWS lambda
data "archive_file" "function_archive_2" {
  type        = "zip"
  source_file = local.gloomahven_companion_service_websocket_binary_path
  output_path = local.gloomahven_companion_service_websocket_archive_path
}

// create the lambda function from zip file
resource "aws_lambda_function" "gloomhaven-companion-service-websocket" {
  function_name = "gloomhaven-companion-service-websocket"
  description   = "Backend websocket service for the gloomhaven companion to support real time scenario updates."
  role          = data.aws_iam_role.lambda_exec.arn
  handler       = local.gloomahven_companion_service_websocket_binary_name
  memory_size   = 128

  filename         = local.gloomahven_companion_service_websocket_archive_path
  source_code_hash = data.archive_file.function_archive_2.output_base64sha256

  runtime = "provided.al2023"
}

resource "aws_cloudwatch_log_group" "gloomhaven-companion-service-websocket" {
  name = "/aws/lambda/${aws_lambda_function.gloomhaven-companion-service-websocket.function_name}"

  retention_in_days = 30
}
