package errors

import (
	"encoding/json"
	"gloomhaven-companion-service/internal/constants"
)

type forbiddenError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e forbiddenError) Error() string {
	jsonString, _ := json.Marshal(e)
	return string(jsonString)
}

func NewForbiddenError() error {
	return forbiddenError{Code: constants.STATUS_CODE_FORBIDDEN, Message: constants.FORBIDDEN_ERROR_MESSAGE}
}
