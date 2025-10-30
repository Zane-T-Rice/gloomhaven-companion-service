// zip the binary, as we can use only zip files to AWS lambda
data "archive_file" "function_archive_3" {
  type        = "zip"
  source_file = local.gloomhaven_companion_service_websocket_default_binary_path
  output_path = local.gloomhaven_companion_service_websocket_default_archive_path
}

// create the lambda function from zip file
resource "aws_lambda_function" "gloomhaven-companion-service-websocket-default" {
  function_name = "gloomhaven-companion-service-websocket-default"
  description   = "Handles broadcasting messages who are connected to the gloomhaven-companion-service-websocket."
  role          = data.aws_iam_role.lambda_exec.arn
  handler       = local.gloomhaven_companion_service_websocket_default_binary_name
  memory_size   = 128

  filename         = local.gloomhaven_companion_service_websocket_default_archive_path
  source_code_hash = data.archive_file.function_archive_3.output_base64sha256

  runtime = "provided.al2023"
}

resource "aws_cloudwatch_log_group" "gloomhaven-companion-service-websocket-default" {
  name = "/aws/lambda/${aws_lambda_function.gloomhaven-companion-service-websocket-default.function_name}"

  retention_in_days = 30
}

# Allow API Gateway to invoke the Lambda
resource "aws_lambda_permission" "allow_api_gateway_default" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.gloomhaven-companion-service-websocket-default.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.websocket_api.execution_arn}/*/*"
}
