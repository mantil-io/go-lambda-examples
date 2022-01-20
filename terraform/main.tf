terraform {
  backend "local" {
    path = "./.state/terraform.tfstate"
  }
}

variable "function_name" {
  type    = string
  default = "go-tf-handler-example"
}

provider "aws" {}

resource "aws_iam_role" "fn" {
  name = "${var.function_name}-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_lambda_function" "fn" {
  role          = aws_iam_role.fn.arn
  function_name = var.function_name
  filename      = "../handler/function.zip"
  runtime       = "provided.al2"
  handler       = "bootstrap"
  architectures = ["arm64"]
}

output "function_name" {
  value = var.function_name
}

output "function_arn" {
  value = aws_lambda_function.fn.arn
}
