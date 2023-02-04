package errs

// PublicHTTPValidationError public Http validation error
type PublicHTTPValidationError struct {
	PublicHTTPError

	// List of errors received while validating payload against schema
	// Required: true
	ValidationErrors []*HTTPValidationErrorDetail `json:"validationErrors"`
}

// HTTPValidationErrorDetail http validation error detail
type HTTPValidationErrorDetail struct {

	// Error describing field validation failure
	// Required: true
	Error *string `json:"error"`

	// Indicates how the invalid field was provided
	// Required: true
	In *string `json:"in"`

	// Key of field failing validation
	// Required: true
	Key *string `json:"key"`
}
