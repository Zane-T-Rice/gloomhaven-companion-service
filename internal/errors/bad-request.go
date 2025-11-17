package errors

import (
	"encoding/json"
	"gloomhaven-companion-service/internal/constants"
)

type BadRequestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e BadRequestError) Error() string {
	jsonString, _ := json.Marshal(e)
	return string(jsonString)
}

func NewBadRequestError(message *string) error {
	errorMessage := constants.BAD_REQUEST_ERROR_MESSAGE
	if message != nil {
		errorMessage = *message
	}
	return BadRequestError{Code: constants.STATUS_CODE_BAD_REQUEST, Message: errorMessage}
}
