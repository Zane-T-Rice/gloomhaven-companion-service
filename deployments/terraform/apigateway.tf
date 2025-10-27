// API Gateway connected to the existing Lambda function
data "aws_caller_identity" "current" {}

resource "aws_api_gateway_rest_api" "gloomhaven_api" {
  name = "gloomhaven-companion-api"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "proxy" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_api.id
  parent_id   = aws_api_gateway_rest_api.gloomhaven_api.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "proxy_any" {
  rest_api_id   = aws_api_gateway_rest_api.gloomhaven_api.id
  resource_id   = aws_api_gateway_resource.proxy.id
  http_method   = "ANY"
  authorization = "NONE"

  request_parameters = {
    "method.request.path.proxy" = true
  }
}

resource "aws_api_gateway_integration" "lambda_proxy" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method.proxy_any.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.gloomhaven-companion-service.invoke_arn
}

resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.gloomhaven-companion-service.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.gloomhaven_api.execution_arn}/*/*"
}

resource "aws_api_gateway_deployment" "gloomhaven_deployment" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_api.id

  # Force new deployment whenever the API definition changes
  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.gloomhaven_api))
  }

  depends_on = [aws_api_gateway_integration.lambda_proxy]
}

resource "aws_api_gateway_stage" "dev" {
  stage_name    = "dev"
  rest_api_id   = aws_api_gateway_rest_api.gloomhaven_api.id
  deployment_id = aws_api_gateway_deployment.gloomhaven_deployment.id
}

output "api_invoke_url" {
  description = "Invoke URL for the API Gateway"
  value       = "https://${aws_api_gateway_rest_api.gloomhaven_api.id}.execute-api.${var.aws_region}.amazonaws.com/${aws_api_gateway_stage.prod.stage_name}"
}
