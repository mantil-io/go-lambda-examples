terraform {
  backend "local" {
    path = "./.state/terraform.tfstate"
  }
}

variable "project_name" {
  type    = string
  default = "api-example"
}

variable "route" {
  type    = string
  default = "handler"
}

locals {
  function_name = "${var.project_name}-${var.route}"
}

provider "aws" {}

output "function_name" {
  value = aws_lambda_function.fn.function_name
}

output "function_arn" {
  value = aws_lambda_function.fn.arn
}

output "endpoint" {
  value = aws_apigatewayv2_api.http.api_endpoint
}

output "url" {
  value = "${aws_apigatewayv2_api.http.api_endpoint}/${var.route}"
}
