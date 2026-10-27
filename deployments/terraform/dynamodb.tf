// DynamoDB table for the gloomhaven companion service
resource "aws_dynamodb_table" "gloomhaven_companion_service" {
  name         = var.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "parent"
  range_key    = "entity"

  // Two independent guards so that undeploying the service can never delete the campaigns stored here.
  //
  // AWS itself refuses to delete the table, whoever asks (terraform, the console or the CLI).
  // Turn it off with deletion_protection_enabled = false and an apply before deleting on purpose.
  deletion_protection_enabled = true

  lifecycle {
    // Terraform refuses any plan that would destroy or replace the table, including terraform destroy.
    // To undeploy everything else, destroy with -target on the other resources or remove the table
    // from the state first (terraform state rm aws_dynamodb_table.gloomhaven_companion_service).
    prevent_destroy = true
  }

  attribute {
    name = "parent"
    type = "S"
  }

  attribute {
    name = "entity"
    type = "S"
  }

  global_secondary_index {
    name            = "entity-index"
    hash_key        = "entity"
    range_key       = "parent"
    projection_type = "KEYS_ONLY" # Or "KEYS_ONLY", "INCLUDE"
  }
}
