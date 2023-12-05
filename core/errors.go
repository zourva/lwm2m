package core

import (
	"errors"
	"github.com/zourva/lwm2m/coap"
)

var (
	ErrorNone                error = nil
	BadRequest                     = errors.New("bad request")
	Unauthorized                   = errors.New("unauthorized")
	BadOption                      = errors.New("bad option")
	Forbidden                      = errors.New("forbidden")
	Conflict                       = errors.New("conflict")
	NotFound                       = errors.New("not found")
	MethodNotAllowed               = errors.New("method not allowed")
	NotAcceptable                  = errors.New("not acceptable")
	RequestEntityIncomplete        = errors.New("request entity incomplete")
	PreconditionFailed             = errors.New("precondition failed")
	RequestEntityTooLarge          = errors.New("request entity too large")
	UnsupportedContentFormat       = errors.New("unsupported content format")
	InternalServerError            = errors.New("internal server error")
	NotImplemented                 = errors.New("not implemented")
	BadGateway                     = errors.New("bad gateway")
	ServiceUnavailable             = errors.New("service unavailable")
	GatewayTimeout                 = errors.New("gateway timeout")
	ProxyingNotSupported           = errors.New("proxying not supported")
)

func GetErrorCode(err error) coap.Code {
	return errorCodesMapping[err]
}

func GetCodeError(code coap.Code) error {
	return codesErrorMapping[code]
}

var errorCodesMapping = map[error]coap.Code{
	ErrorNone:                coap.CodeEmpty,
	BadRequest:               coap.CodeBadRequest,
	Unauthorized:             coap.CodeUnauthorized,
	BadOption:                coap.CodeBadOption,
	Forbidden:                coap.CodeForbidden,
	Conflict:                 coap.CodeConflict,
	NotFound:                 coap.CodeNotFound,
	MethodNotAllowed:         coap.CodeMethodNotAllowed,
	NotAcceptable:            coap.CodeNotAcceptable,
	RequestEntityIncomplete:  coap.CodeRequestEntityIncomplete,
	PreconditionFailed:       coap.CodePreconditionFailed,
	RequestEntityTooLarge:    coap.CodeRequestEntityTooLarge,
	UnsupportedContentFormat: coap.CodeUnsupportedMediaType,
	InternalServerError:      coap.CodeInternalServerError,
	NotImplemented:           coap.CodeNotImplemented,
	BadGateway:               coap.CodeBadGateway,
	ServiceUnavailable:       coap.CodeServiceUnavailable,
	GatewayTimeout:           coap.CodeGatewayTimeout,
	ProxyingNotSupported:     coap.CodeProxyingNotSupported,
}

var codesErrorMapping = map[coap.Code]error{
	coap.CodeEmpty:                   ErrorNone,
	coap.CodeBadRequest:              BadRequest,
	coap.CodeUnauthorized:            Unauthorized,
	coap.CodeBadOption:               BadOption,
	coap.CodeForbidden:               Forbidden,
	coap.CodeConflict:                Conflict,
	coap.CodeNotFound:                NotFound,
	coap.CodeMethodNotAllowed:        MethodNotAllowed,
	coap.CodeNotAcceptable:           NotAcceptable,
	coap.CodeRequestEntityIncomplete: RequestEntityIncomplete,
	coap.CodePreconditionFailed:      PreconditionFailed,
	coap.CodeRequestEntityTooLarge:   RequestEntityTooLarge,
	coap.CodeUnsupportedMediaType:    UnsupportedContentFormat,
	coap.CodeInternalServerError:     InternalServerError,
	coap.CodeNotImplemented:          NotImplemented,
	coap.CodeBadGateway:              BadGateway,
	coap.CodeServiceUnavailable:      ServiceUnavailable,
	coap.CodeGatewayTimeout:          GatewayTimeout,
	coap.CodeProxyingNotSupported:    ProxyingNotSupported,
}
