# define new API Gateway
# ref: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_api
resource "aws_apigatewayv2_api" "http" {
  name          = "${var.project_name}-http"
  protocol_type = "HTTP"
  cors_configuration {
    allow_origins = toset(["*"])
  }
}

# create CloudWatch group for the API Gateway access logs
# ref: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group
resource "aws_cloudwatch_log_group" "access_logs" {
  name              = "/aws/vendedlogs/${aws_apigatewayv2_api.http.name}-access-logs"
  retention_in_days = 14
}

# create API Gateway stage: https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-stages.html
# ref: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_stage
resource "aws_apigatewayv2_stage" "default" {
  name        = "$default"
  api_id      = aws_apigatewayv2_api.http.id
  auto_deploy = true
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.access_logs.arn
    # using default format
    # ref: https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-logging-variables.html
    format = jsonencode({
      httpMethod     = "$context.httpMethod"
      ip             = "$context.identity.sourceIp"
      protocol       = "$context.protocol"
      requestId      = "$context.requestId"
      requestTime    = "$context.requestTime"
      responseLength = "$context.responseLength"
      routeKey       = "$context.routeKey"
      status         = "$context.status"
    })
  }
}

# connect API Gateway to the Lambda function
# ref: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_integration
resource "aws_apigatewayv2_integration" "proxy" {
  api_id                 = aws_apigatewayv2_api.http.id
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  payload_format_version = "2.0"
  integration_uri        = aws_lambda_function.fn.invoke_arn
}

# expose Lambda function at some route
# ref: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_route
resource "aws_apigatewayv2_route" "route" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "ANY /${var.route}"
  target    = "integrations/${aws_apigatewayv2_integration.proxy.id}"
}

# give API Gateway permission to call our Lambda function
#ref: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_permission
resource "aws_lambda_permission" "api_gateway_invoke" {
  function_name = aws_lambda_function.fn.function_name
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http.execution_arn}/*/*"
}
