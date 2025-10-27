// DynamoDB table for the gloomhaven companion service
resource "aws_dynamodb_table" "gloomhaven_companion" {
  name         = var.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"
  range_key    = "entity"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "entity"
    type = "S"
  }
}
