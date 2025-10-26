scripts/build.sh
terraform -chdir=deployments/terraform init
terraform -chdir=deployments/terraform apply