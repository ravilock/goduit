package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheck(t *testing.T) {
	e := echo.New()
	t.Run("Should return healthcheck", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := Healthcheck(c)
		assert.NoError(t, err)
		if rec.Code != http.StatusOK {
			t.Errorf("Got status different than %q, got %q", http.StatusOK, rec.Code)
		}
		expectedBody := "OK\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("Expected body to equal %q, got %q", expectedBody, rec.Body.String())
		}
	})
}
