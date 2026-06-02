output "api_endpoint" {
  description = "The Live URL for the Cope and Hope API"
  value       = "${aws_apigatewayv2_api.http_api.api_endpoint}/weather?city=London&mode=cope"
}
