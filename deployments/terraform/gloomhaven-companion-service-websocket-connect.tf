// create the lambda function from zip file
resource "aws_lambda_function" "gloomhaven_companion_service_websocket_connect" {
  function_name = "gloomhaven-companion-service-websocket-connect"
  description   = "Handles connecting new clients to the gloomhaven-companion-service-websocket."
  role          = data.aws_iam_role.lambda_exec.arn
  handler       = local.gloomhaven_companion_service_websocket_connect_binary_name
  memory_size   = 128

  filename         = local.gloomhaven_companion_service_websocket_connect_archive_path
  source_code_hash = data.archive_file.function_archive_2.output_base64sha256

  runtime = "provided.al2023"
}

resource "aws_cloudwatch_log_group" "gloomhaven_companion_service_websocket_connect" {
  name = "/aws/lambda/${aws_lambda_function.gloomhaven_companion_service_websocket_connect.function_name}"

  retention_in_days = 30
}

# Allow API Gateway to invoke the Lambda
resource "aws_lambda_permission" "allow_api_gateway_connect" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.gloomhaven_companion_service_websocket_connect.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.websocket_api.execution_arn}/*/*"
}
