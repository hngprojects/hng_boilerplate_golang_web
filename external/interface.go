package external

type SendExternalRequest interface {
	// SetupExternalRequest prepares the external request struct with necessary fields
	SetupExternalRequest(name string, data any) error
	// SendRequest sends the request to the external services
	SendRequest() (any, error)
}
