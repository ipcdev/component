package errors

import (
	"fmt"
	"net/http"
	"sync"
)

// codes contains a map of error codes to metadata.
var (
	codeMux = &sync.Mutex{}
	codes   = map[int]Coder{}
)

var (
	unknownCoder = defaultCoder{
		HTTP: http.StatusInternalServerError,
		C:    0,
		Ext:  "An internal server error occurred",
	}
)

func init() {
	codes[unknownCoder.Code()] = unknownCoder
}

// Coder defines an interface for an error code detail information.
type Coder interface {
	// HTTPStatus that should be used for the associated error code.
	HTTPStatus() int

	// Code returns the code of the coder
	Code() int

	// String external (user) facing error text.
	String() string
}

type defaultCoder struct {
	// HTTP status that should be used for the associated error code.
	HTTP int

	// C refers to the integer code of the ErrCode.
	C int

	// External (user) facing error text.
	Ext string
}

// HTTPStatus returns the associated HTTP status code, if any. Otherwise,
// returns 200.
func (coder defaultCoder) HTTPStatus() int {
	if coder.HTTP == 0 {
		return 500
	}

	return coder.HTTP
}

// Code returns the integer code of the coder.
func (coder defaultCoder) Code() int {
	return coder.C
}

// String implements stringer. String returns the external error message,
// if any.
func (coder defaultCoder) String() string {
	return coder.Ext
}

// Register a user define error code.
// It will override to exist code.
func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved as unknownCode error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

// MustRegister register a user define error code.
// It will panic when the same Code already exist.
func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is reserved as unknownCode error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

// ParseCoder parse any error into *withCode.
// nil error will return nil direct.
// None withStack error will be parsed as ErrUnknown.
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}

	return unknownCoder
}

// IsCode reports whether any error in err's chain contains the given error code.
func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, code)
		}

		return false
	}

	return false
}
