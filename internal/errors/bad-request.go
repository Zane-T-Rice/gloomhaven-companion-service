package errors

import (
	"encoding/json"
	"gloomhaven-companion-service/internal/constants"
)

type badRequestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e badRequestError) Error() string {
	jsonString, _ := json.Marshal(e)
	return string(jsonString)
}

func NewBadRequestError() error {
	return badRequestError{Code: constants.STATUS_CODE_BAD_REQUEST, Message: constants.BAD_REQUEST_ERROR_MESSAGE}
}
