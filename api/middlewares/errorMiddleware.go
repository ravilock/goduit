package middlewares

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

const defaultStatusCode = http.StatusInternalServerError
const defaultErrorMessage = "Internal Server Error"

func ErrorMiddleware(err error, c echo.Context) {
	c.Logger().Error(err)

	if httpError := new(echo.HTTPError); errors.As(err, &httpError) {
		if httpError.Code < http.StatusInternalServerError {
			c.String(httpError.Code, httpError.Error())
			return
		} else {
			httpError.Internal = nil
			c.String(httpError.Code, httpError.Error())
			return
		}
	}

	if mongo.IsDuplicateKeyError(err) {
		c.String(http.StatusConflict, "Content Already Exists")
		return
	}

	c.String(defaultStatusCode, defaultErrorMessage)
}
