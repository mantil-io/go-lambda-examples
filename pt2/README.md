
cd handler
../../scripts/build.sh

cd ../terraform
terraform init

export AWS_REGION=eu-central-1
export AWS_PROFILE=org5

terraform apply
