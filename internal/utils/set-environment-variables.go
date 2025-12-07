package utils

import (
	"github.com/joho/godotenv"
)

type Secret struct {
	Secret string `json:"value"`
}

func SetEnvironmentVariables() {
	godotenv.Load()
}
