package integrationtests

import (
	"strings"

	"github.com/google/uuid"
)

func UniqueUsername(args ...string) string {
	if len(args) == 0 {
		return uuid.NewString()
	}
	return uuid.NewString() + strings.Join(args, "")
}

func UniqueEmail(args ...string) string {
	prefix := UniqueUsername(args...)
	return prefix + "@test.test"
}

func UniqueTitle(args ...string) string {
	base := UniqueEmail(args...)
	return strings.Join(strings.Split(base, "-"), " ")
}
