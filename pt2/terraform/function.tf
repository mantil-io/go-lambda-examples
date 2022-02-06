# define Lmabda execution role
# A Lambda function's execution role is an AWS Identity and Access Management (IAM)
# role that grants the function permission to access AWS services and resources.
# You provide this role when you create a function, and Lambda assumes the role when your function is invoked.
# ref: https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html
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

# add policy to the IAM role
# ref: https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html
resource "aws_iam_role_policy_attachment" "basic" {
  role       = aws_iam_role.fn.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# create CloudWatch log group for the function
resource "aws_cloudwatch_log_group" "function_log_group" {
  name              = "/aws/lambda/${local.function_name}"
  retention_in_days = 14
}

# define our function
resource "aws_lambda_function" "fn" {
  role             = aws_iam_role.fn.arn
  function_name    = local.function_name
  filename         = "../handler/function.zip"
  source_code_hash = filebase64sha256("../handler/function.zip")
  runtime          = "provided.al2"
  handler          = "bootstrap"
  architectures    = ["arm64"]
}
