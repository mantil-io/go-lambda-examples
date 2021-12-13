# AWS cli

In the `handler`` folder is dummy Lambda function handler. That is copy of the code from AWS [docs](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html). 

Shell scripts in the same folder are examples of how to use AWS cli.

publish.sh will build, prepare deployment package and then create new Lambda function or if already exists update code existing. 

cleanup.sh will delete IAM role and Lambda created in publish

invoke.sh is example of how to run this new Lambda function. Payload parameter is input to the function and output is in response.json file.
This function is not exposed throught API Gateway, doesn't have any integrations so the only way to invoke it is through console or cli/sdk call. 


handler/main.go is code from AWS Lambda function handler [example](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html)


# Terraform

./build.sh

terraform init
terraform apply --auto-approve

../handler/invoke.sh  $(terraform output --raw function_name)

terraform destroy --auto-approve