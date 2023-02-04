package errs

import (
	"net/http"

	"github.com/driif/echo-go-starter/internal/server/net/errs"
)

var (
	UserExists   = errs.NewHTTPError(http.StatusConflict, "USER_ALREADY_EXISTS", "User already exists.")
	UserNotFound = errs.NewHTTPError(http.StatusNotFound, "USER_NOT_FOUND", "User not found.")
)
