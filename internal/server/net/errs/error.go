package errs

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// HTTPErrorTypeGeneric is a generic error type returned as default for all HTTP errors without a specific type.
	HTTPErrorTypeGeneric = "generic"
)

// Payload is accordance with RFC 7807 (Problem Details for HTTP APIs) with the exception of the type value not being represented as a URI.
// https://tools.ietf.org/html/rfc7807

type HTTPError struct {
	PublicHTTPError
	Internal       error                  `json:"-"`
	AdditionalData map[string]interface{} `json:"-"`
}

type HTTPValidationError struct {
	PublicHTTPValidationError
	Internal       error                  `json:"-"`
	AdditionalData map[string]interface{} `json:"-"`
}

// NewHTTPError creates a new HTTPError with the given code, type and title.
func NewHTTPError(code int, errorType, title string) *HTTPError {
	return &HTTPError{
		PublicHTTPError: PublicHTTPError{
			Code:  &code,
			Type:  &errorType,
			Title: &title,
		},
	}
}

// NewHTTPErrorWithDetail creates a new HTTPError with the given code, type, title and detail.
func NewHTTPErrorWithDetail(code int, errorType, title, detail string) *HTTPError {
	return &HTTPError{
		PublicHTTPError: PublicHTTPError{
			Code:   &code,
			Type:   &errorType,
			Title:  &title,
			Detail: detail,
		},
	}
}

// Returns the error message from HTTPError.
func (e *HTTPError) Error() string {
	var b strings.Builder

	fmt.Fprintf(&b, "HTTPError %d (%s): %s", *e.Code, *e.Type, *e.Title)

	if len(e.Detail) > 0 {
		fmt.Fprintf(&b, " - %s", e.Detail)
	}
	if e.Internal != nil {
		fmt.Fprintf(&b, ", %v", e.Internal)
	}
	if e.AdditionalData != nil && len(e.AdditionalData) > 0 {
		keys := make([]string, 0, len(e.AdditionalData))
		for k := range e.AdditionalData {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteString(". Additional: ")
		for i, k := range keys {
			fmt.Fprintf(&b, "%s=%v", k, e.AdditionalData[k])
			if i < len(keys)-1 {
				b.WriteString(", ")
			}
		}
	}

	return b.String()
}

// NewHTTPValidationError creates a new HTTPValidationError with the given code, type, title and validation errors.
func NewHTTPValidationError(code int, errorType, title string, validationErrors []*HTTPValidationErrorDetail) *HTTPValidationError {
	return &HTTPValidationError{
		PublicHTTPValidationError: PublicHTTPValidationError{
			PublicHTTPError: PublicHTTPError{
				Code:  &code,
				Type:  &errorType,
				Title: &title,
			},
			ValidationErrors: validationErrors,
		},
	}
}

// NewHTTPValidationErrorWithDetail creates a new HTTPValidationError with the given code, type, title, detail and validation errors.
func NewHTTPValidationErrorWithDetail(code int, errorType, title, detail string, validationErrors []*HTTPValidationErrorDetail) *HTTPValidationError {
	return &HTTPValidationError{
		PublicHTTPValidationError: PublicHTTPValidationError{
			PublicHTTPError: PublicHTTPError{
				Code:   &code,
				Type:   &errorType,
				Title:  &title,
				Detail: detail,
			},
			ValidationErrors: validationErrors,
		},
	}
}

// Returns the error message from HTTPValidationError.
func (e *HTTPValidationError) Error() string {
	var b strings.Builder

	fmt.Fprintf(&b, "HTTPValidationError %d (%s): %s", *e.Code, *e.Type, *e.Title)

	if len(e.Detail) > 0 {
		fmt.Fprintf(&b, " - %s", e.Detail)
	}
	if e.Internal != nil {
		fmt.Fprintf(&b, ", %v", e.Internal)
	}
	if e.AdditionalData != nil && len(e.AdditionalData) > 0 {
		keys := make([]string, 0, len(e.AdditionalData))
		for k := range e.AdditionalData {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		b.WriteString(". Additional: ")
		for i, k := range keys {
			fmt.Fprintf(&b, "%s=%v", k, e.AdditionalData[k])
			if i < len(keys)-1 {
				b.WriteString(", ")
			}
		}
	}

	b.WriteString(" - Validation: ")
	for i, ve := range e.ValidationErrors {
		fmt.Fprintf(&b, "%s (in %s): %s", *ve.Key, *ve.In, *ve.Error)
		if i < len(e.ValidationErrors)-1 {
			b.WriteString(", ")
		}
	}

	return b.String()
}
