package openweather

import (
	"fmt"
)

// ErrAPIError is returned when the API returns an error
type ErrAPIError struct {
	Err  error
	Msg  string
	Code int
}

func (e *ErrAPIError) Error() string {
	if e.Msg == "" {
		e.Msg = "api error"
	}
	if e.Code != 0 {
		e.Msg += fmt.Sprintf(" (%d)", e.Code)
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrNoAPIKey is returned when no API key is provided
type ErrNoAPIKey struct {
	Err error
	Msg string
}

// Error returns the error message
func (e *ErrNoAPIKey) Error() string {
	if e.Msg == "" {
		e.Msg = "no api key provided- use WithAPIKey()"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}

// ErrNoLocation is returned when no location is provided
type ErrNoLocation struct {
	Err error
	Msg string
}

// Error returns the error message
func (e *ErrNoLocation) Error() string {
	if e.Msg == "" {
		e.Msg = "no location provided- use WithLocation()"
	}
	if e.Err != nil {
		e.Msg += ": " + e.Err.Error()
	}
	return e.Msg
}
