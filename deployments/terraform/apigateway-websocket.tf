// This is the websocket API Gateway.
// This connects integrations for the $connect, $disconnect, and $default routes.

resource "aws_apigatewayv2_api" "websocket_api" {
  name                       = "gloomhaven-companion-service-websocket"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_integration" "connect_integration" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.gloomhaven_companion_service_websocket_connect.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_integration" "default_integration" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.gloomhaven_companion_service_websocket_default.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_integration" "disconnect_integration" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.gloomhaven_companion_service_websocket_disconnect.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "connect" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  route_key          = "$connect"
  authorization_type = "NONE"
  target             = "integrations/${aws_apigatewayv2_integration.connect_integration.id}"
}

resource "aws_apigatewayv2_route" "default" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  route_key          = "$default"
  authorization_type = "NONE"
  target             = "integrations/${aws_apigatewayv2_integration.default_integration.id}"
}

resource "aws_apigatewayv2_route" "disconnect" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  route_key          = "$disconnect"
  authorization_type = "NONE"
  target             = "integrations/${aws_apigatewayv2_integration.disconnect_integration.id}"
}

resource "aws_apigatewayv2_deployment" "websocket_deployment" {
  api_id = aws_apigatewayv2_api.websocket_api.id
  triggers = {
    redeploy = sha1(jsonencode({
      routes       = [aws_apigatewayv2_route.connect.id, aws_apigatewayv2_route.default.id, aws_apigatewayv2_route.disconnect.id],
      integrations = [aws_apigatewayv2_integration.connect_integration.id, aws_apigatewayv2_integration.default_integration.id, aws_apigatewayv2_integration.disconnect_integration.id]
    }))
  }
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [aws_apigatewayv2_route.connect, aws_apigatewayv2_route.default, aws_apigatewayv2_route.disconnect, aws_apigatewayv2_integration.connect_integration, aws_apigatewayv2_integration.default_integration, aws_apigatewayv2_integration.disconnect_integration]
}