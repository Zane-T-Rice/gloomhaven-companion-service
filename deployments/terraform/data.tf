/* 
  This role needs to be manually created.
  This role needs the following policies:
    "logs:CreateLogGroup"
    "logs:CreateLogStream"
    "logs:PutLogEvents"
    "lambda:InvokeFunction"
    "dynamodb:PutItem"
    "dynamodb:DeleteItem"
    "dynamodb:GetItem"
    "dynamodb:Query",
    "dynamodb:UpdateItem"
    "execute-api:ManageConnections"
    "secretsmanager:GetSecretValue"
*/
data "aws_iam_role" "lambda_exec" {
  name = "AWSLambdaBasicExecutionRole"
}
// These are local file paths for the go binaries created by scripts/build.sh.
data "archive_file" "function_archive" {
  type        = "zip"
  source_file = local.gloomhaven_companion_service_binary_path
  output_path = local.gloomhaven_companion_service_archive_path
}
data "archive_file" "function_archive_2" {
  type        = "zip"
  source_file = local.gloomhaven_companion_service_websocket_connect_binary_path
  output_path = local.gloomhaven_companion_service_websocket_connect_archive_path
}
data "archive_file" "function_archive_3" {
  type        = "zip"
  source_file = local.gloomhaven_companion_service_websocket_default_binary_path
  output_path = local.gloomhaven_companion_service_websocket_default_archive_path
}
data "archive_file" "function_archive_4" {
  type        = "zip"
  source_file = local.gloomhaven_companion_service_websocket_disconnect_binary_path
  output_path = local.gloomhaven_companion_service_websocket_disconnect_archive_path
}
