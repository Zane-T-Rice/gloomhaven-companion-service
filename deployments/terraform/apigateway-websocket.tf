// WebSocket API Gateway (APIGWv2) wired to existing Lambdas.
// - Uses `aws_lambda_function.gloomhaven-companion-service-websocket` as a REQUEST (custom) authorizer
// - Integrates routes to the main `gloomhaven-companion-service` Lambda via AWS_PROXY

resource "aws_apigatewayv2_api" "websocket_api" {
  name                       = "gloomhaven-companion-service-websocket"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_integration" "connect_integration" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.gloomhaven-companion-service-websocket-connect.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_integration" "default_integration" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.gloomhaven-companion-service-websocket-default.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_integration" "disconnect_integration" {
  api_id             = aws_apigatewayv2_api.websocket_api.id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.gloomhaven-companion-service-websocket-disconnect.invoke_arn
  integration_method = "POST"
}


// resource "aws_apigatewayv2_authorizer" "lambda_request_authorizer" {
//   api_id                            = aws_apigatewayv2_api.websocket_api.id
//   authorizer_type                   = "REQUEST"
//   name                              = "websocket-lambda-request-authorizer"
//   authorizer_uri                    = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/${aws_lambda_function.gloomhaven-companion-service-websocket-connect.arn}/invocations"
//   identity_sources                  = ["route.request.header.Authorization"]
// }

resource "aws_apigatewayv2_route" "connect" {
  api_id    = aws_apigatewayv2_api.websocket_api.id
  route_key = "$connect"

  authorization_type = "NONE"
  // authorizer_id      = aws_apigatewayv2_authorizer.lambda_request_authorizer.id

  target = "integrations/${aws_apigatewayv2_integration.connect_integration.id}"
}

resource "aws_apigatewayv2_route" "default" {
  api_id    = aws_apigatewayv2_api.websocket_api.id
  route_key = "$default"

  authorization_type = "NONE"

  target = "integrations/${aws_apigatewayv2_integration.default_integration.id}"
}

resource "aws_apigatewayv2_route" "disconnect" {
  api_id    = aws_apigatewayv2_api.websocket_api.id
  route_key = "$disconnect"

  authorization_type = "NONE"

  target = "integrations/${aws_apigatewayv2_integration.disconnect_integration.id}"
}



// resource "aws_lambda_permission" "allow_apigw_invoke_main" {
// 	statement_id  = "AllowExecutionFromWebSocketAPIGWMain"
// 	action        = "lambda:InvokeFunction"
// 	function_name = aws_lambda_function.gloomhaven-companion-service.function_name
// 	principal     = "apigateway.amazonaws.com"
// 	source_arn    = "${aws_apigatewayv2_api.websocket_api.execution_arn}/*/*"
// }
// 
// resource "aws_lambda_permission" "allow_apigw_invoke_authorizer" {
// 	statement_id  = "AllowExecutionFromWebSocketAPIGWAuthorizer"
// 	action        = "lambda:InvokeFunction"
// 	function_name = aws_lambda_function.gloomhaven-companion-service-websocket.function_name
// 	principal     = "apigateway.amazonaws.com"
// 	source_arn    = "${aws_apigatewayv2_api.websocket_api.execution_arn}/*/*"
// }

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

resource "aws_apigatewayv2_stage" "prod" {
  api_id        = aws_apigatewayv2_api.websocket_api.id
  name          = "prod"
  deployment_id = aws_apigatewayv2_deployment.websocket_deployment.id
}
