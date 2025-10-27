// API Gateway for the gloomhaven-companion-service Lambda
resource "aws_api_gateway_rest_api" "gloomhaven_api" {
  name        = "gloomhaven-companion-api"
  description = "REST API for gloomhaven companion service backed by Lambda"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

# Root ANY method to proxy the root path /
resource "aws_api_gateway_method" "root_any" {
  rest_api_id   = aws_api_gateway_rest_api.gloomhaven_api.id
  resource_id   = aws_api_gateway_rest_api.gloomhaven_api.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "root_lambda" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_api.id
  resource_id = aws_api_gateway_rest_api.gloomhaven_api.root_resource_id
  http_method = aws_api_gateway_method.root_any.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.gloomhaven-companion-service.invoke_arn
}

# Proxy resource to route all nested paths
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

resource "aws_api_gateway_integration" "proxy_lambda" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_api.id
  resource_id = aws_api_gateway_resource.proxy.id
  http_method = aws_api_gateway_method.proxy_any.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.gloomhaven-companion-service.invoke_arn
}

# Allow API Gateway to invoke the Lambda
resource "aws_lambda_permission" "allow_api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.gloomhaven-companion-service.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.gloomhaven_api.execution_arn}/*/*"
}

# Deployment and stage
resource "aws_api_gateway_deployment" "gloomhaven_deployment" {
  rest_api_id = aws_api_gateway_rest_api.gloomhaven_api.id

  # redeploy when the API definition changes
  triggers = {
    redeploy = sha1(jsonencode({
      resources = aws_api_gateway_rest_api.gloomhaven_api
    }))
  }

  depends_on = [aws_api_gateway_integration.proxy_lambda, aws_api_gateway_integration.root_lambda]
}

resource "aws_api_gateway_stage" "prod" {
  stage_name    = "prod"
  rest_api_id   = aws_api_gateway_rest_api.gloomhaven_api.id
  deployment_id = aws_api_gateway_deployment.gloomhaven_deployment.id

  # optional: enable access logging or settings here
}

output "api_rest_api_id" {
  description = "API Gateway REST API id"
  value       = aws_api_gateway_rest_api.gloomhaven_api.id
}

output "api_invoke_url" {
  description = "Invoke URL for the API Gateway (execute-api)"
  value       = "https://${aws_api_gateway_rest_api.gloomhaven_api.id}.execute-api.${var.aws_region}.amazonaws.com/${aws_api_gateway_stage.prod.stage_name}"
}
