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

output "cert_validation_cname_name" {
  description = "The CNAME name to add to your DNS for validation"
  value       = tolist(aws_acm_certificate.cert.domain_validation_options)[0].resource_record_name
}

output "cert_validation_cname_value" {
  description = "The CNAME value to add to your DNS for validation"
  value       = tolist(aws_acm_certificate.cert.domain_validation_options)[0].resource_record_value
}

output "custom_domain_target" {
  description = "Point demo.kagno.com to this URL with a CNAME record AFTER deployment succeeds"
  value       = aws_apigatewayv2_domain_name.custom_domain.domain_name_configuration[0].target_domain_name
}
