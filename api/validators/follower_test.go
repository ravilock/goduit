package validators

import (
	"testing"

	"github.com/ravilock/goduit/api"
	"github.com/ravilock/goduit/api/requests"
)

func TestFollower(t *testing.T) {
	InitValidator()
	t.Run("Username is required", func(t *testing.T) {
		request := &requests.Follower{}

		got := Follower(request)
		assertError(t, got, api.RequiredFieldError("Username"))
	})
	t.Run("Username must not be blank", func(t *testing.T) {
		request := &requests.Follower{Username: "   "}

		got := Follower(request)
		assertError(t, got, api.RequiredFieldError("Username"))
	})
	t.Run("Username must have at least 5 chars", func(t *testing.T) {
		request := &requests.Follower{Username: "1234"}

		got := Follower(request)
		assertError(t, got, api.InvalidFieldLength("Username", "min", "5"))
	})
	t.Run("Username must have at most 255 chars", func(t *testing.T) {
		request := &requests.Follower{Username: randomString(256)}

		got := Follower(request)
		assertError(t, got, api.InvalidFieldLength("Username", "max", "255"))
	})
}
