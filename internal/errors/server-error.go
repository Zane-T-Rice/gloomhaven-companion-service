package errors

import (
	"encoding/json"
	"gloomhaven-companion-service/internal/constants"
)

type serverError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e serverError) Error() string {
	jsonString, _ := json.Marshal(e)
	return string(jsonString)
}

func NewServerError() error {
	return serverError{Code: constants.STATUS_CODE_SERVER_ERROR, Message: constants.SERVER_ERROR_ERROR_MESSAGE}
}
