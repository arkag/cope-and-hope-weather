resource "aws_apigatewayv2_api" "http_api" {
  name          = "CopeAndHopeGateway"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "$default"
  auto_deploy = true
}
