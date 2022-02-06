resource "aws_apigatewayv2_api" "http" {
  name          = "${var.project_name}-http"
  protocol_type = "HTTP"
  cors_configuration {
    allow_origins = toset(["*"])
  }
}

resource "aws_cloudwatch_log_group" "http_access_logs" {
  name              = "/aws/vendedlogs/${aws_apigatewayv2_api.http.name}-access-logs"
  retention_in_days = 14
}

resource "aws_apigatewayv2_stage" "http_default" {
  name          = "$default"
  api_id        = aws_apigatewayv2_api.http.id
  deployment_id = aws_apigatewayv2_deployment.http.id

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.http_access_logs.arn
    format          = "{ \"requestId\":\"$context.requestId\", \"ip\": \"$context.identity.sourceIp\", \"requestTime\":\"$context.requestTime\", \"httpMethod\":\"$context.httpMethod\",\"routeKey\":\"$context.routeKey\", \"status\":\"$context.status\",\"protocol\":\"$context.protocol\", \"responseLength\":\"$context.responseLength\" }"
  }

  default_route_settings {
    detailed_metrics_enabled = true
    throttling_burst_limit   = 5000
    throttling_rate_limit    = 10000
  }
}

resource "aws_apigatewayv2_route" "http_proxy" {
  api_id    = aws_apigatewayv2_api.http.id
  route_key = "ANY /${var.route}/{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.http_proxy.id}"
}

resource "aws_apigatewayv2_integration" "http_proxy" {
  api_id             = aws_apigatewayv2_api.http.id
  integration_type   = "AWS_PROXY" #each.value.type
  integration_method = "POST"      #each.value.integration_method
  integration_uri    = aws_lambda_function.fn.invoke_arn
  #TODO: ???
  request_parameters = {
    "overwrite:path" = "$request.path.proxy"
  }
}

resource "aws_apigatewayv2_deployment" "http" {
  api_id = aws_apigatewayv2_api.http.id
  depends_on = [
    aws_apigatewayv2_integration.http_proxy,
    aws_apigatewayv2_route.http_proxy,
  ]
  triggers = {
    redeployment = sha1(join(",", tolist([
      jsonencode(aws_apigatewayv2_integration.http_proxy),
      jsonencode(aws_apigatewayv2_route.http_proxy),
    ])))
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lambda_permission" "api_gateway_invoke" {
  function_name = aws_lambda_function.fn.function_name
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http.execution_arn}/*/*"
}
