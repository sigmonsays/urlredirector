package main

import (
	"errors"
	"fmt"
)

//go:generate stringer -type=ErrorClass

type ErrorClass int

const (
	Unknown        ErrorClass = -1
	Success        ErrorClass = 0
	NoSuchRedirect ErrorClass = 40001
)

var (
	UnknownError = &ApiError{
		Err:   errors.New("Unknown"),
		Class: Unknown,
	}
)

func (me ErrorClass) Errorf(e error, message string, args ...interface{}) *ApiError {
	ret := NewApiError(e)
	ret.Class = me
	ret.Message = fmt.Sprintf(message, args...)
	return ret
}
func (me ErrorClass) Sprintf(message string, args ...interface{}) *ApiError {
	return me.Errorf(errors.New("ERROR"), message, args...)
}

func (me ErrorClass) Int32() int32 {
	return int32(me)
}

type ApiError struct {
	// the actual error that occurred
	Err error

	// the class of error
	Class ErrorClass

	Message string
}

func (me *ApiError) Error() string {
	return fmt.Sprintf("ApiError(Classs:%s/%d message:%q err:%q)", me.Class, me.Class, me.Message, me.Err)
}

func IsApiError(e error) *ApiError {
	ie, ok := e.(*ApiError)
	if ok {
		return ie
	}
	return UnknownError
}

// holds an underlying error and necessary information
func NewApiError(e error) *ApiError {
	ret := &ApiError{
		Err:   e,
		Class: Unknown,
	}
	return ret
}
