data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda_exec" {
  name               = "cope_hope_lambda_exec"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

# Attach basic execution policy so Lambda can write logs to CloudWatch
resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "dynamo_access" {
  statement {
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem"
    ]
    # Restrict the Lambda so it can ONLY touch our specific cache table
    resources = [aws_dynamodb_table.weather_cache.arn]
  }
}

resource "aws_iam_role_policy" "dynamo_access" {
  name   = "dynamo_access"
  role   = aws_iam_role.lambda_exec.id
  policy = data.aws_iam_policy_document.dynamo_access.json
}
