// API Gateway for the gloomhaven-companion-service Lambda
// This proxies all calls through to the Lambda.
resource "aws_api_gateway_rest_api" "gloomhaven_companion_service" {
  name                         = "gloomhaven-companion-service"
  description                  = "REST API for gloomhaven companion service backed by Lambda"
  disable_execute_api_endpoint = true

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}
resource "aws_api_gateway_method" "root_any" {
  rest_api_id   = aws_api_gateway_rest_api.gloomhaven_companion_service.id
  resource_id   = aws_api_gateway_rest_api.gloomhaven_companion_service.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
}
resource "aws_api_gateway_integration" "root_lambda" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_companion_service.id
  resource_id = aws_api_gateway_rest_api.gloomhaven_companion_service.root_resource_id
  http_method = aws_api_gateway_method.root_any.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.gloomhaven_companion_service.invoke_arn
}
resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_companion_service.id
  parent_id   = aws_api_gateway_rest_api.gloomhaven_companion_service.root_resource_id
  path_part   = "{proxy+}"
}
resource "aws_api_gateway_method" "proxy_any" {
  rest_api_id   = aws_api_gateway_rest_api.gloomhaven_companion_service.id
  resource_id   = aws_api_gateway_resource.proxy.id
  http_method   = "ANY"
  authorization = "NONE"

  request_parameters = {
    "method.request.path.proxy" = true
  }
}
resource "aws_api_gateway_integration" "proxy_lambda" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_companion_service.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method.proxy_any.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.gloomhaven_companion_service.invoke_arn
}
resource "aws_lambda_permission" "allow_api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.gloomhaven_companion_service.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.gloomhaven_companion_service.execution_arn}/*/*"
}
resource "aws_api_gateway_deployment" "gloomhaven_deployment" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_companion_service.id
  # redeploy when the API definition changes
  triggers = {
    redeploy = sha1(jsonencode({
      resources = aws_api_gateway_rest_api.gloomhaven_companion_service
    }))
  }
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [aws_api_gateway_integration.proxy_lambda, aws_api_gateway_integration.root_lambda]
}