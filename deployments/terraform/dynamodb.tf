// DynamoDB table for the gloomhaven companion service
resource "aws_dynamodb_table" "gloomhaven-companion-service" {
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

  global_secondary_index {
    name            = "entity-index"
    hash_key        = "entity"
    range_key       = "id"
    projection_type = "KEYS_ONLY" # Or "KEYS_ONLY", "INCLUDE"
  }
}
