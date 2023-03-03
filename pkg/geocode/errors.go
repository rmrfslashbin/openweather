package geocode

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
