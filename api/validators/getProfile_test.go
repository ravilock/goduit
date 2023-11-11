package validators

import (
	"math/rand"
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestGetProfile(t *testing.T) {
	t.Run("Username is required", func(t *testing.T) {
		request := &requests.GetProfile{}

		got := GetProfile(request)
		assertError(t, got, api.RequiredFieldError("Username"))
	})
	t.Run("Username must not be blank", func(t *testing.T) {
		request := &requests.GetProfile{Username: "   "}

		got := GetProfile(request)
		assertError(t, got, api.RequiredFieldError("Username"))
	})
	t.Run("Username must have at least 5 chars", func(t *testing.T) {
		request := &requests.GetProfile{Username: "1234"}

		got := GetProfile(request)
		assertError(t, got, api.InvalidFieldLength("Username", "min", "5"))
	})
	t.Run("Username must have at most 255 chars", func(t *testing.T) {
		request := &requests.GetProfile{Username: randomString(256)}

		got := GetProfile(request)
		assertError(t, got, api.InvalidFieldLength("Username", "max", "255"))
	})
	t.Run("Valid username should not return error", func(t *testing.T) {
		request := &requests.GetProfile{Username: randomString(40)}

		got := GetProfile(request)
		assertNoError(t, got)
	})
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got == nil {
		t.Errorf("Expected error but didn't receive one")
		return
	}

	if got.Error() != want.Error() {
		t.Errorf("Got %q want %q", got, want)
		return
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()

	if got != nil {
		t.Errorf("Got an error but did'nt want one")
	}
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
