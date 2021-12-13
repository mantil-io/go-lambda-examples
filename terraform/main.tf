terraform {
  backend "local" {
    path = "./.state/terraform.tfstate"
  }
}

locals {
  aws_region    = "eu-central-1"
  function_name = "go-tf-handler-example"
}

provider "aws" {
  region = local.aws_region
}

resource "aws_iam_role" "fn" {
  name = "${local.function_name}-role"

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
  function_name = local.function_name
  filename      = "../handler/function.zip"
  runtime       = "provided.al2"
  handler       = "bootstrap"
  architectures = ["arm64"]
}


output "function_name" {
  value = local.function_name
}

output "function_arn" {
  value = aws_lambda_function.fn.arn
}
