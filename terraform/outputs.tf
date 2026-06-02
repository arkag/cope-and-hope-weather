output "api_gateway_id" {
  description = "The API Gateway ID for component terraform to reference"
  value       = aws_apigatewayv2_api.http_api.id
}

output "api_gateway_execution_arn" {
  description = "The API Gateway execution ARN for Lambda permissions"
  value       = aws_apigatewayv2_api.http_api.execution_arn
}

output "api_endpoint" {
  description = "The public base URL"
  value       = aws_apigatewayv2_api.http_api.api_endpoint
}
