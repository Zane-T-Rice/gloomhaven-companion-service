GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ../dist/bootstrap ../cmd/gloomhaven-companion-service/
terraform -chdir=../deployments/terraform init
terraform -chdir=../deployments/terraform apply