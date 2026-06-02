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

resource "aws_iam_role" "web_lambda_exec" {
  name               = "cope_hope_web_lambda_exec"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_role_policy_attachment" "web_lambda_basic_execution" {
  role       = aws_iam_role.web_lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "invoke_api" {
  statement {
    actions   = ["lambda:InvokeFunction"]
    resources = ["arn:aws:lambda:us-east-1:*:function:CopeAndHopeWeatherAPI"]
  }
}

resource "aws_iam_role_policy" "invoke_api" {
  name   = "invoke_api"
  role   = aws_iam_role.web_lambda_exec.id
  policy = data.aws_iam_policy_document.invoke_api.json
}

# --- Lambda Function ---
data "archive_file" "dummy_web_lambda" {
  type        = "zip"
  output_path = "${path.module}/dummy.zip"
  source {
    content  = "dummy payload to initialize terraform"
    filename = "bootstrap"
  }
}

resource "aws_lambda_function" "web" {
  function_name    = "CopeAndHopeWebUI"
  role             = aws_iam_role.web_lambda_exec.arn
  handler          = "app.handler"
  runtime          = "python3.12"
  filename         = data.archive_file.dummy_web_lambda.output_path
  source_code_hash = data.archive_file.dummy_web_lambda.output_base64sha256

  environment {
    variables = {
      API_ENDPOINT = data.aws_apigatewayv2_api.http_api.api_endpoint
    }
  }

  lifecycle {
    ignore_changes = [filename, source_code_hash]
  }
}

resource "aws_cloudwatch_log_group" "web_logs" {
  name              = "/aws/lambda/${aws_lambda_function.web.function_name}"
  retention_in_days = 7
}

# --- API Gateway Integration ---
resource "aws_apigatewayv2_integration" "web_integration" {
  api_id             = local.api_gateway_id
  integration_type   = "AWS_PROXY"
  integration_uri    = aws_lambda_function.web.invoke_arn
  integration_method = "POST"
}

# Root route serves the web UI
resource "aws_apigatewayv2_route" "root_route" {
  api_id    = local.api_gateway_id
  route_key = "GET /"
  target    = "integrations/${aws_apigatewayv2_integration.web_integration.id}"
}

# Search endpoint (proxied by Flask to the Go API)
resource "aws_apigatewayv2_route" "search_route" {
  api_id    = local.api_gateway_id
  route_key = "GET /search"
  target    = "integrations/${aws_apigatewayv2_integration.web_integration.id}"
}

# Health check
resource "aws_apigatewayv2_route" "health_route" {
  api_id    = local.api_gateway_id
  route_key = "GET /health"
  target    = "integrations/${aws_apigatewayv2_integration.web_integration.id}"
}

# Static assets (CSS, JS)
resource "aws_apigatewayv2_route" "static_route" {
  api_id    = local.api_gateway_id
  route_key = "GET /static/{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.web_integration.id}"
}

resource "aws_lambda_permission" "web_gw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.web.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${data.aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}
