package setenvironmentvariables

import (
	"context"
	"encoding/json"
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
	log.Println("IN SetEnvironmentVariables")
	godotenv.Load()
	localServicePort := os.Getenv("LOCAL_SERVICE_PORT")
	if localServicePort == "" {
		config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
		if err != nil {
			log.Fatal(err)
		}

		secretManagerClient := secretsmanager.NewFromConfig(config)

		setEnvironmentVariable("gloomhaven-companion-service-audience", "AUDIENCE", secretManagerClient)
		setEnvironmentVariable("gloomhaven-companion-service-issuer", "ISSUER", secretManagerClient)
		setEnvironmentVariable("gloomhaven-companion-service-url", "GLOOMHAVEN_COMPANION_SERVICE_URL", secretManagerClient)
	}
}

func setEnvironmentVariable(secretName string, environmentVariableName string, svc *secretsmanager.Client) {
	log.Println("IN setEnvironmentVariable")
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
