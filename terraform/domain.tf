variable "domain_name" {
  description = "The custom domain name"
  default     = "demo.kagno.com"
}

# 1. Request the Certificate
resource "aws_acm_certificate" "cert" {
  domain_name       = var.domain_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

# 2. Validate the Certificate (This pauses Terraform!)
resource "aws_acm_certificate_validation" "cert_validation" {
  certificate_arn         = aws_acm_certificate.cert.arn
  validation_record_fqdns = [for record in aws_acm_certificate.cert.domain_validation_options : record.resource_record_name]
}

# 3. Create the API Gateway Custom Domain
resource "aws_apigatewayv2_domain_name" "custom_domain" {
  domain_name = var.domain_name

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.cert_validation.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

# 4. Map the Custom Domain to the API Gateway
resource "aws_apigatewayv2_api_mapping" "api_mapping" {
  api_id      = aws_apigatewayv2_api.http_api.id
  domain_name = aws_apigatewayv2_domain_name.custom_domain.id
  stage       = aws_apigatewayv2_stage.default.id
}
