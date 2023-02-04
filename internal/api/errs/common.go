package errs

import (
	"net/http"

	"github.com/driif/echo-go-starter/internal/server/net/errs"
)

var (
	NotUUID   = errs.NewHTTPError(http.StatusExpectationFailed, "NOT_UUID", "Not a valid UUID.")
	ParseBody = errs.NewHTTPError(http.StatusExpectationFailed, "PARSE_BODY", "Could not parse body.")
)
