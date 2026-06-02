# Create a dummy zip file. Terraform requires a payload to create a Lambda function,
# but we want GitHub Actions to deploy the *real* code later.
data "archive_file" "dummy_lambda" {
  type        = "zip"
  output_path = "${path.module}/dummy.zip"
  source {
    content  = "dummy payload to initialize terraform"
    filename = "bootstrap"
  }
}

resource "aws_lambda_function" "api" {
  function_name    = "CopeAndHopeWeatherAPI"
  role             = aws_iam_role.lambda_exec.arn
  handler          = "bootstrap"   # The compiled Go binary must be named 'bootstrap' in Amazon Linux 2023
  runtime          = "provided.al2023" 
  filename         = data.archive_file.dummy_lambda.output_path
  source_code_hash = data.archive_file.dummy_lambda.output_base64sha256

  environment {
    variables = {
      OWM_API_KEY = var.owm_api_key
    }
  }

  lifecycle {
    # Ignore changes to the code zip. This prevents Terraform from overwriting 
    # the live application code every time you run `terraform apply` for infrastructure changes.
    ignore_changes = [filename, source_code_hash]
  }
}

# Explicitly create the log group so we can set a 7-day retention policy (Free Tier safe)
resource "aws_cloudwatch_log_group" "api_logs" {
  name              = "/aws/lambda/${aws_lambda_function.api.function_name}"
  retention_in_days = 7 
}
