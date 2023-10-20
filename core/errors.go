package core

import "errors"

type ErrorType = string

const (
	ClientNotFound      ErrorType = "ClientNotFound"
	OperationNotAllowed ErrorType = "OperationNotAllowed"

	ErrorNone                ErrorType = ""
	BadRequest               ErrorType = "BadRequest"
	Unauthorized             ErrorType = "Unauthorized"
	BadOption                ErrorType = "BadOption"
	Forbidden                ErrorType = "Forbidden"
	Conflict                 ErrorType = "Conflict"
	NotFound                 ErrorType = "NotFound"
	MethodNotAllowed         ErrorType = "MethodNotAllowed"
	NotAcceptable            ErrorType = "NotAcceptable"
	RequestEntityIncomplete  ErrorType = "RequestEntityIncomplete"
	PreconditionFailed       ErrorType = "PreconditionFailed"
	RequestEntityTooLarge    ErrorType = "RequestEntityTooLarge"
	UnsupportedContentFormat ErrorType = "UnsupportedContentFormat"
)

var errorMap = map[ErrorType]error{
	ClientNotFound:      errors.New("registered client not found"),
	OperationNotAllowed: errors.New("operation not allowed"),
}

func Errors(t ErrorType) error {
	return errorMap[t]
}
