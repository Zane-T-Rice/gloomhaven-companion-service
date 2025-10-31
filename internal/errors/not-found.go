package errors

import (
	"encoding/json"
	"gloomhaven-companion-service/internal/constants"
)

type notFoundError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e notFoundError) Error() string {
	jsonString, _ := json.Marshal(e)
	return string(jsonString)
}

func NewNotFoundError() error {
	return notFoundError{Code: constants.STATUS_CODE_NOT_FOUND, Message: constants.NOT_FOUND_ERROR_MESSAGE}
}
