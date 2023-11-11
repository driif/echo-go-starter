package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/driif/echo-go-starter/internal/server/net/errs"
	"github.com/driif/echo-go-starter/pkg/logs"
	"github.com/driif/echo-go-starter/pkg/strs"
	"github.com/labstack/echo/v4"
)

var (
	// DefaultHTTPErrorHandlerConfig is the default config for the HTTPErrorHandler
	DefaultHTTPErrorHandlerConfig = HTTPErrorHandlerConfig{
		HideInternalServerErrorDetails: false,
	}
)

// HTTPErrorHandlerConfig is the config for the HTTPErrorHandler
type HTTPErrorHandlerConfig struct {
	HideInternalServerErrorDetails bool
}

// HTTPErrorHandler is a custom HTTP error handler
func HTTPErrorHandler() echo.HTTPErrorHandler {
	return HTTPErrorHandlerWithConfig(DefaultHTTPErrorHandlerConfig)
}

// HTTPErrorHandlerWithConfig is a custom HTTP error handler with config
func HTTPErrorHandlerWithConfig(config HTTPErrorHandlerConfig) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var code int
		var he error

		var httpError *errs.HTTPError
		var httpValidationError *errs.HTTPValidationError
		var echoHTTPError *echo.HTTPError

		switch {
		case errors.As(err, &httpError):
			code = *httpError.Code
			he = httpError

			if code == http.StatusInternalServerError && config.HideInternalServerErrorDetails {
				if httpError.Internal == nil {
					//nolint:errorlint
					httpError.Internal = fmt.Errorf("internal error: %s", httpError)
				}
				title := http.StatusText(http.StatusInternalServerError)
				httpError.Title = &title
			}
		case errors.As(err, &httpValidationError):
			code = *httpValidationError.Code
			he = httpValidationError

			if code == http.StatusInternalServerError && config.HideInternalServerErrorDetails {
				if httpValidationError.Internal == nil {
					//nolint:errorlint
					httpValidationError.Internal = fmt.Errorf("internal error: %s", httpValidationError)
				}

				title := http.StatusText(http.StatusInternalServerError)
				httpValidationError.Title = &title
			}
		case errors.As(err, &echoHTTPError):
			code = echoHTTPError.Code

			if code == http.StatusInternalServerError && config.HideInternalServerErrorDetails {
				if echoHTTPError.Internal == nil {
					//nolint:errorlint
					echoHTTPError.Internal = fmt.Errorf("internal error: %s", echoHTTPError)
				}

				he = &errs.HTTPError{
					PublicHTTPError: errs.PublicHTTPError{
						Code:  &echoHTTPError.Code,
						Title: strs.StrToPtr(http.StatusText(http.StatusInternalServerError)),
						Type:  strs.StrToPtr(errs.HTTPErrorTypeGeneric),
					},
					Internal: echoHTTPError.Internal,
				}
			} else {
				msg, ok := echoHTTPError.Message.(string)
				if !ok {
					if m, errr := json.Marshal(msg); errr == nil {
						msg = string(m)
					} else {
						msg = fmt.Sprintf("failed to marshal HTTP error message: %v", errr)
					}
				}

				he = &errs.HTTPError{
					PublicHTTPError: errs.PublicHTTPError{
						Code:  &echoHTTPError.Code,
						Title: &msg,
						Type:  strs.StrToPtr(errs.HTTPErrorTypeGeneric),
					},
					Internal: echoHTTPError.Internal,
				}
			}
		default:
			code = http.StatusInternalServerError

			if config.HideInternalServerErrorDetails {
				he = &errs.HTTPError{
					PublicHTTPError: errs.PublicHTTPError{
						Code:  &code,
						Title: strs.StrToPtr(http.StatusText(http.StatusInternalServerError)),
						Type:  strs.StrToPtr(errs.HTTPErrorTypeGeneric),
					},
					Internal: err,
				}
			} else {
				he = &errs.HTTPError{
					PublicHTTPError: errs.PublicHTTPError{
						Code:  &code,
						Title: strs.StrToPtr(err.Error()),
						Type:  strs.StrToPtr(errs.HTTPErrorTypeGeneric),
					},
				}
			}
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, he)
			}

			if err != nil {
				logs.LogFromEchoContext(c).Warn().Err(err).AnErr("http_err", err).Msg("Failed to handle HTTP error")
			}
		}
	}
}
