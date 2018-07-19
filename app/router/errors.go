package router

import (
	"fmt"
	"net/http"

	"github.com/epointpayment/mloc-cpe/app/codes"
	"github.com/epointpayment/mloc-cpe/app/config"
	"github.com/epointpayment/mloc-cpe/app/controllers"
	"github.com/epointpayment/mloc-cpe/app/log"

	"github.com/juju/errors"
	"github.com/labstack/echo"
)

var (
	ErrUnsupportedMediaType = codes.New("ERROR_UNSUPPORTED_MEDIA_TYPE").
				WithMessage(http.StatusText(http.StatusUnsupportedMediaType)).
				WithStatusCode(http.StatusUnsupportedMediaType).
				RegisterError()

	ErrNotFound = codes.New("ERROR_NOT_FOUND").
			WithMessage(http.StatusText(http.StatusNotFound)).
			WithStatusCode(http.StatusNotFound).
			RegisterError()

	ErrUnauthorized = codes.New("ERROR_UNAUTHORIZED").
			WithMessage(http.StatusText(http.StatusUnauthorized)).
			WithStatusCode(http.StatusUnauthorized).
			RegisterError()

	ErrForbidden = codes.New("ERROR_FORBIDDEN").
			WithMessage(http.StatusText(http.StatusForbidden)).
			WithStatusCode(http.StatusForbidden).
			RegisterError()

	ErrMethodNotAllowed = codes.New("ERROR_METHOD_NOT_ALLOWED").
				WithMessage(http.StatusText(http.StatusMethodNotAllowed)).
				WithStatusCode(http.StatusMethodNotAllowed).
				RegisterError()

	ErrStatusRequestEntityTooLarge = codes.New("ERROR_REQUEST_ENTITY_TOO_LARGE").
					WithMessage(http.StatusText(http.StatusRequestEntityTooLarge)).
					WithStatusCode(http.StatusRequestEntityTooLarge).
					RegisterError()

	ErrUnknown = codes.New("ERROR_UNKNOWN").
			WithMessage("Unknown Error").
			RegisterError()
)

// appendErrorHandler handles errors for the router
func (r *Router) appendErrorHandler() {
	r.e.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{}
			res  codes.Code
		)

		if he, ok := err.(*echo.HTTPError); ok {
			switch he {
			case echo.ErrUnsupportedMediaType:
				err = errors.Wrap(err, ErrUnsupportedMediaType)

			case echo.ErrNotFound:
				err = errors.Wrap(err, ErrNotFound)

			case echo.ErrUnauthorized:
				err = errors.Wrap(err, ErrUnauthorized)

			case echo.ErrForbidden:
				err = errors.Wrap(err, ErrForbidden)

			case echo.ErrMethodNotAllowed:
				err = errors.Wrap(err, ErrMethodNotAllowed)

			case echo.ErrStatusRequestEntityTooLarge:
				err = errors.Wrap(err, ErrStatusRequestEntityTooLarge)

			default:
				code = he.Code
				msg = he.Message
				if he.Internal != nil {
					msg = fmt.Sprintf("%v, %v", err, he.Internal)
				}
			}
		}

		if ae, ok := errors.Cause(err).(codes.Code); ok {
			code = ae.StatusCode
			res = ae
			if !c.Response().Committed {
				log.Debug(err.(*errors.Err).Underlying())
			}
		} else {
			res = ErrUnknown
		}

		if !c.Response().Committed {
			if res == ErrUnknown {
				if msg != nil {
				} else if r.e.Debug {
					msg = err.Error()
				} else {
					msg = http.StatusText(code)
				}
				log.Error(msg)
			}

			if res == ErrUnknown || config.IsDev() {
				if output := log.StackTrace(err); len(output) > 0 {
					log.Print("Stack Trace:\n" + output.String())
				}
			}

			// Send response
			if c.Request().Method == echo.HEAD { // Issue #608
				err = c.NoContent(code)
			} else {
				// Send error in a specific format
				err = controllers.SendErrorResponse(c, res)
			}
			if err != nil {
				r.e.Logger.Error(err)
			}
		}
	}
}
