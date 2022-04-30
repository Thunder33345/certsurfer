package certsurfer

import (
	"errors"
	"fmt"
)

func IsDialError(err error) bool {
	return errors.As(err, &dialError{})
}

type dialError struct {
	err error
}

func (e dialError) Error() string {
	return fmt.Sprintf("error dailing: %v", e.err)
}

func (e dialError) Unwrap() error {
	return e.err
}

func IsReadError(err error) bool {
	return errors.As(err, &readError{})
}

type readError struct {
	err error
}

func (e readError) Error() string {
	return fmt.Sprintf("error reading: %v", e.err)
}

func (e readError) Unwrap() error {
	return e.err
}

func IsJsonError(err error) bool {
	return errors.As(err, &jsonError{})
}

type jsonError struct {
	err error
}

func (e jsonError) Error() string {
	return fmt.Sprintf("error unmarshalling: %v", e.err)
}

func (e jsonError) Unwrap() error {
	return e.err
}
