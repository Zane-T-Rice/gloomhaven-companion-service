// If you do not want to map a custom domain to the api
// and websocket gateways then you can delete or edit this file.

// These custom domains will have to be made manually in 
// API Gateway -> Custom domain names.
data "aws_api_gateway_domain_name" "domain_name" {
  domain_name = local.gloomhaven_companion_service_api_domain_name
}
data "aws_api_gateway_domain_name" "ws_domain_name" {
  domain_name = local.gloomhaven_companion_service_websocket_domain_name
}
// Map custom domain to the API Gateway.
resource "aws_api_gateway_stage" "prod" {
  stage_name    = "prod"
  rest_api_id   = aws_api_gateway_rest_api.gloomhaven_companion_service.id
  deployment_id = aws_api_gateway_deployment.gloomhaven_deployment.id
}
resource "aws_api_gateway_base_path_mapping" "gloomhaven_companion_service" {
  api_id      = aws_api_gateway_rest_api.gloomhaven_companion_service.id
  stage_name  = aws_api_gateway_stage.prod.stage_name
  domain_name = data.aws_api_gateway_domain_name.domain_name.domain_name
  base_path   = aws_lambda_function.gloomhaven_companion_service.function_name
}
// Map custom domain to the web socket gateway.
resource "aws_apigatewayv2_stage" "prod" {
  api_id        = aws_apigatewayv2_api.websocket_api.id
  name          = "prod"
  deployment_id = aws_apigatewayv2_deployment.websocket_deployment.id
}
resource "aws_apigatewayv2_api_mapping" "gloomhaven_companion_service_ws" {
  api_id          = aws_apigatewayv2_api.websocket_api.id
  stage           = aws_apigatewayv2_stage.prod.name
  domain_name     = data.aws_api_gateway_domain_name.ws_domain_name.domain_name
  api_mapping_key = aws_lambda_function.gloomhaven_companion_service.function_name
}
