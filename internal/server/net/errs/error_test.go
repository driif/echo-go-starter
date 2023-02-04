package errs_test

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/driif/echo-go-starter/internal/server/net/errs"
	"github.com/driif/echo-go-starter/pkg/strs"
	"github.com/stretchr/testify/require"
)

func TestHTTPErrorSimple(t *testing.T) {
	e := errs.NewHTTPError(http.StatusNotFound, errs.HTTPErrorTypeGeneric, http.StatusText(http.StatusNotFound))
	require.Equal(t, "HTTPError 404 (generic): Not Found", e.Error())
}

func TestHTTPErrorDetail(t *testing.T) {
	e := errs.NewHTTPErrorWithDetail(http.StatusNotFound, errs.HTTPErrorTypeGeneric, http.StatusText(http.StatusNotFound), "ToS violation")
	require.Equal(t, "HTTPError 404 (generic): Not Found - ToS violation", e.Error())
}

func TestHTTPErrorInternalError(t *testing.T) {
	e := errs.NewHTTPError(http.StatusInternalServerError, errs.HTTPErrorTypeGeneric, http.StatusText(http.StatusInternalServerError))

	e.Internal = sql.ErrConnDone

	require.Equal(t, "HTTPError 500 (generic): Internal Server Error, sql: connection is already closed", e.Error())
}

func TestHTTPErrorAdditionalData(t *testing.T) {
	e := errs.NewHTTPError(http.StatusInternalServerError, errs.HTTPErrorTypeGeneric, http.StatusText(http.StatusInternalServerError))

	e.AdditionalData = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	require.Equal(t, "HTTPError 500 (generic): Internal Server Error. Additional: key1=value1, key2=value2", e.Error())
}

var valErrs = append(make([]*errs.HTTPValidationErrorDetail, 0, 2), &errs.HTTPValidationErrorDetail{
	Key:   strs.StrToPtr("test1"),
	In:    strs.StrToPtr("body.test1"),
	Error: strs.StrToPtr("ValidationError"),
}, &errs.HTTPValidationErrorDetail{
	Key:   strs.StrToPtr("test2"),
	In:    strs.StrToPtr("body.test2"),
	Error: strs.StrToPtr("Validation Error"),
})

func TestHTTPValidationErrorSimple(t *testing.T) {
	e := errs.NewHTTPValidationError(http.StatusBadRequest, errs.HTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)
	require.Equal(t, "HTTPValidationError 400 (generic): Bad Request - Validation: test1 (in body.test1): ValidationError, test2 (in body.test2): Validation Error", e.Error())
}

func TestHTTPValidationErrorDetail(t *testing.T) {
	e := errs.NewHTTPValidationErrorWithDetail(http.StatusBadRequest, errs.HTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), "Did API spec change?", valErrs)
	require.Equal(t, "HTTPValidationError 400 (generic): Bad Request - Did API spec change? - Validation: test1 (in body.test1): ValidationError, test2 (in body.test2): Validation Error", e.Error())
}

func TestHTTPValidationErrorInternalError(t *testing.T) {
	e := errs.NewHTTPValidationError(http.StatusBadRequest, errs.HTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)

	e.Internal = sql.ErrConnDone

	require.Equal(t, "HTTPValidationError 400 (generic): Bad Request, sql: connection is already closed - Validation: test1 (in body.test1): ValidationError, test2 (in body.test2): Validation Error", e.Error())
}

func TestHTTPValidationErrorAdditionalData(t *testing.T) {
	e := errs.NewHTTPValidationError(http.StatusBadRequest, errs.HTTPErrorTypeGeneric, http.StatusText(http.StatusBadRequest), valErrs)

	e.AdditionalData = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	require.Equal(t, "HTTPValidationError 400 (generic): Bad Request. Additional: key1=value1, key2=value2 - Validation: test1 (in body.test1): ValidationError, test2 (in body.test2): Validation Error", e.Error())
}
