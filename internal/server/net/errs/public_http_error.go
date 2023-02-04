package errs

// PublicHTTPError public Http error
type PublicHTTPError struct {

	// HTTP status code returned for the error
	// Example: 403
	// Required: true
	// Maximum: 599
	// Minimum: 100
	Code *int `json:"status"`

	// More detailed, human-readable, optional explanation of the error
	// Example: User is lacking permission to access this resource
	Detail string `json:"detail,omitempty"`

	// Short, human-readable description of the error
	// Example: Forbidden
	// Required: true
	Title *string `json:"title"`

	// Type of error returned, should be used for client-side error handling
	// Example: generic
	// Required: true
	Type *string `json:"type"`
}
