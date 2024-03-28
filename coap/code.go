package coap

import "strconv"

// Code defines a valid CoAP code which
// is comprised by Method and Response codes.
type Code uint8

// 3-bits method codes;
// 5-bits response codes
const (
	Get    Code = 1
	Post   Code = 2
	Put    Code = 3
	Delete Code = 4

	// 2.xx

	CodeEmpty    Code = 0  // 2.00
	CodeCreated  Code = 65 // 2.01
	CodeDeleted  Code = 66 // 2.02
	CodeValid    Code = 67 // 2.03
	CodeChanged  Code = 68 // 2.04
	CodeContent  Code = 69 // 2.05
	CodeContinue Code = 95 // 2.31

	// 4.xx

	CodeBadRequest              Code = 128 // 4.00
	CodeUnauthorized            Code = 129 // 4.01
	CodeBadOption               Code = 130 // 4.02
	CodeForbidden               Code = 131 // 4.03
	CodeNotFound                Code = 132 // 4.04
	CodeMethodNotAllowed        Code = 133 // 4.05
	CodeNotAcceptable           Code = 134 // 4.06
	CodeRequestEntityIncomplete Code = 136 // 4.08
	CodeConflict                Code = 137 // 4.09
	CodePreconditionFailed      Code = 140 // 4.12
	CodeRequestEntityTooLarge   Code = 141 // 4.13
	CodeUnsupportedMediaType    Code = 143 // 4.15
	CodeTooManyRequests         Code = 157 // 4.29

	// 5.xx

	CodeInternalServerError  Code = 160 // 5.00
	CodeNotImplemented       Code = 161 // 5.01
	CodeBadGateway           Code = 162 // 5.02
	CodeServiceUnavailable   Code = 163 // 5.03
	CodeGatewayTimeout       Code = 164 // 5.04
	CodeProxyingNotSupported Code = 165 // 5.05
)

var codeToString = map[Code]string{
	Get:                       "GET",
	Post:                      "POST",
	Put:                       "PUT",
	Delete:                    "DELETE",
	CodeEmpty:                 "0 Empty",
	CodeCreated:               "201 Created",
	CodeDeleted:               "202 Deleted",
	CodeValid:                 "203 Valid",
	CodeChanged:               "204 Changed",
	CodeContent:               "205 Content",
	CodeBadRequest:            "400 BadRequest",
	CodeUnauthorized:          "401 Unauthorized",
	CodeBadOption:             "402 BadOption",
	CodeForbidden:             "403 Forbidden",
	CodeNotFound:              "404 NotFound",
	CodeMethodNotAllowed:      "405 MethodNotAllowed",
	CodeNotAcceptable:         "406 NotAcceptable",
	CodePreconditionFailed:    "412 PreconditionFailed",
	CodeRequestEntityTooLarge: "413 RequestEntityTooLarge",
	CodeUnsupportedMediaType:  "415 UnsupportedMediaType",
	CodeTooManyRequests:       "429 TooManyRequests",
	CodeInternalServerError:   "500 InternalServerError",
	CodeNotImplemented:        "501 NotImplemented",
	CodeBadGateway:            "502 BadGateway",
	CodeServiceUnavailable:    "503 ServiceUnavailable",
	CodeGatewayTimeout:        "504 GatewayTimeout",
	CodeProxyingNotSupported:  "505 ProxyingNotSupported",
}

// CodeString returns the string representation of a Code
func CodeString(c Code) string {
	return c.String()
}

func (c Code) String() string {
	val, ok := codeToString[c]
	if ok {
		return val
	}
	return "Code(" + strconv.FormatInt(int64(c), 10) + ")"
}

// Created returns true if coap code is Created.
func (c Code) Created() bool                 { return c == CodeCreated }
func (c Code) Deleted() bool                 { return c == CodeDeleted }
func (c Code) Valid() bool                   { return c == CodeValid }
func (c Code) Changed() bool                 { return c == CodeChanged }
func (c Code) Content() bool                 { return c == CodeContent }
func (c Code) Continue() bool                { return c == CodeContinue }
func (c Code) BadRequest() bool              { return c == CodeBadRequest }
func (c Code) Unauthorized() bool            { return c == CodeUnauthorized }
func (c Code) BadOption() bool               { return c == CodeBadOption }
func (c Code) Forbidden() bool               { return c == CodeForbidden }
func (c Code) NotFound() bool                { return c == CodeNotFound }
func (c Code) MethodNotAllowed() bool        { return c == CodeMethodNotAllowed }
func (c Code) NotAcceptable() bool           { return c == CodeNotAcceptable }
func (c Code) RequestEntityIncomplete() bool { return c == CodeRequestEntityIncomplete }
func (c Code) PreconditionFailed() bool      { return c == CodePreconditionFailed }
func (c Code) RequestEntityTooLarge() bool   { return c == CodeRequestEntityTooLarge }
func (c Code) UnsupportedMediaType() bool    { return c == CodeUnsupportedMediaType }
func (c Code) TooManyRequests() bool         { return c == CodeTooManyRequests }
func (c Code) InternalServerError() bool     { return c == CodeInternalServerError }
func (c Code) NotImplemented() bool          { return c == CodeNotImplemented }
func (c Code) BadGateway() bool              { return c == CodeBadGateway }
func (c Code) ServiceUnavailable() bool      { return c == CodeServiceUnavailable }
func (c Code) GatewayTimeout() bool          { return c == CodeGatewayTimeout }
func (c Code) ProxyingNotSupported() bool    { return c == CodeProxyingNotSupported }
