package handler

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		max       int
		wantValue string
		wantErr   bool
	}{
		{"accepts a normal value", "Backend Engineer", 160, "Backend Engineer", false},
		{"trims surrounding space", "  Acme Inc  ", 160, "Acme Inc", false},
		{"rejects empty", "", 160, "", true},
		{"rejects whitespace only", "   \t\n ", 160, "", true},
		{"rejects over the limit", strings.Repeat("x", 161), 160, strings.Repeat("x", 161), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fe := fieldErrors{}
			got := fe.required("title", tc.input, tc.max)
			assert.Equal(t, tc.wantValue, got)
			assert.Equal(t, tc.wantErr, !fe.ok())
		})
	}
}

func TestEmail(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantValue string
		wantErr   bool
	}{
		{"accepts a valid address", "alice@example.com", "alice@example.com", false},
		{"lowercases so logins match", "Alice@Example.COM", "alice@example.com", false},
		{"rejects a missing domain", "alice@", "alice@", true},
		{"rejects text without an at sign", "not-an-email", "not-an-email", true},
		{"rejects empty", "", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fe := fieldErrors{}
			got := fe.email("email", tc.input, 320)
			assert.Equal(t, tc.wantValue, got)
			assert.Equal(t, tc.wantErr, !fe.ok())
		})
	}
}

func TestOneOf(t *testing.T) {
	t.Run("defaults when empty", func(t *testing.T) {
		fe := fieldErrors{}
		assert.Equal(t, "open", fe.oneOf("status", "", "open", "open", "paused", "closed"))
		assert.True(t, fe.ok())
	})

	t.Run("accepts an allowed value", func(t *testing.T) {
		fe := fieldErrors{}
		assert.Equal(t, "closed", fe.oneOf("status", "closed", "open", "open", "paused", "closed"))
		assert.True(t, fe.ok())
	})

	t.Run("rejects a value outside the set", func(t *testing.T) {
		fe := fieldErrors{}
		fe.oneOf("status", "banana", "open", "open", "paused", "closed")
		assert.False(t, fe.ok())
		assert.Contains(t, fe, "status")
	})
}

func TestFieldErrorsCollectsEveryProblem(t *testing.T) {
	fe := fieldErrors{}
	fe.required("title", "", 160)
	fe.required("company", "   ", 160)
	fe.email("email", "nope", 320)

	assert.False(t, fe.ok())
	assert.Len(t, fe, 3)
}
