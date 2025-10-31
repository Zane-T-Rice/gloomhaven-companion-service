package utils

import (
	"context"
	"encoding/json"
	"gloomhaven-companion-service/internal/constants"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
)

type Secret struct {
	Secret string `json:"value"`
}

func SetEnvironmentVariables() {
	godotenv.Load()
	localServicePort := os.Getenv("LOCAL_SERVICE_PORT")
	if localServicePort == "" {
		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
		if err != nil {
			log.Fatal(err)
		}

		secretManagerClient := secretsmanager.NewFromConfig(config)

		setEnvironmentVariable(constants.AUDIENCE_SECRET_NAME, constants.AUDIENCE, secretManagerClient)
		setEnvironmentVariable(constants.ISSUER_SECRET_NAME, constants.ISSUER, secretManagerClient)
		setEnvironmentVariable(constants.URL_SECRET_NAME, constants.GLOOMHAVEN_COMPANION_SERVICE_URL, secretManagerClient)
		setEnvironmentVariable(constants.WEBSITE_DOMAIN_SECRET_NAME, constants.WEBSITE_DOMAIN, secretManagerClient)
		setEnvironmentVariable(constants.API_GATEWAY_BASE_PATH_SECRET_NAME, constants.API_GATEWAY_BASE_PATH, secretManagerClient)
	}
}

func setEnvironmentVariable(secretName string, environmentVariableName string, svc *secretsmanager.Client) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	secret := Secret{}
	json.Unmarshal([]byte(*result.SecretString), &secret)
	os.Setenv(environmentVariableName, secret.Secret)
}
