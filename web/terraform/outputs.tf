output "web_url" {
  description = "The live URL for the Cope and Hope Web UI"
  value       = "${data.aws_apigatewayv2_api.http_api.api_endpoint}/"
}

output "health_url" {
  description = "The health check endpoint"
  value       = "${data.aws_apigatewayv2_api.http_api.api_endpoint}/health"
}
