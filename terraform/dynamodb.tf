resource "aws_dynamodb_table" "weather_cache" {
  name           = "WeatherCache"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "city"

  attribute {
    name = "city"
    type = "S"
  }

  ttl {
    attribute_name = "expires_at"
    enabled        = true
  }
}
