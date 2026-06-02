# --- Look up the shared API Gateway ---
data "aws_apigatewayv2_apis" "shared" {
  protocol_type = "HTTP"
  name          = "CopeAndHopeGateway"
}

locals {
  api_gateway_id = tolist(data.aws_apigatewayv2_apis.shared.ids)[0]
}

data "aws_apigatewayv2_api" "http_api" {
  api_id = local.api_gateway_id
}

# --- DynamoDB Cache Table ---
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

# --- IAM Role ---
data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "api_lambda_exec" {
  name               = "cope_hope_api_lambda_exec"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_role_policy_attachment" "api_lambda_basic_execution" {
  role       = aws_iam_role.api_lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "dynamo_access" {
  statement {
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem"
    ]
    resources = [aws_dynamodb_table.weather_cache.arn]
  }
}

resource "aws_iam_role_policy" "dynamo_access" {
  name   = "dynamo_access"
  role   = aws_iam_role.api_lambda_exec.id
  policy = data.aws_iam_policy_document.dynamo_access.json
}

# --- Lambda Function ---
data "archive_file" "dummy_api_lambda" {
  type        = "zip"
  output_path = "${path.module}/dummy.zip"
  source {
    content  = "dummy payload to initialize terraform"
    filename = "bootstrap"
  }
}

resource "aws_lambda_function" "api" {
  function_name    = "CopeAndHopeWeatherAPI"
  role             = aws_iam_role.api_lambda_exec.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  filename         = data.archive_file.dummy_api_lambda.output_path
  source_code_hash = data.archive_file.dummy_api_lambda.output_base64sha256

  environment {
    variables = {
      OWM_API_KEY = var.owm_api_key
    }
  }

  lifecycle {
    ignore_changes = [filename, source_code_hash]
  }
}

resource "aws_cloudwatch_log_group" "api_logs" {
  name              = "/aws/lambda/${aws_lambda_function.api.function_name}"
  retention_in_days = 7
}

# --- API Gateway Integration ---
resource "aws_apigatewayv2_integration" "api_integration" {
  api_id             = local.api_gateway_id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.api.invoke_arn
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "weather_route" {
  api_id    = local.api_gateway_id
  route_key = "GET /weather"
  target    = "integrations/${aws_apigatewayv2_integration.api_integration.id}"
}

resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${data.aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}
