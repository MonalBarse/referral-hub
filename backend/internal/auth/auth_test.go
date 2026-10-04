package auth

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueAndParseRoundTrip(t *testing.T) {
	m := New("test-secret")
	userID := uuid.New()

	token, err := m.Issue(userID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	got, err := m.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, userID, got)
}

func TestParseRejectsForeignSignature(t *testing.T) {
	issuer := New("the-real-secret")
	verifier := New("a-different-secret")

	token, err := issuer.Issue(uuid.New())
	require.NoError(t, err)

	_, err = verifier.Parse(token)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestParseRejectsGarbage(t *testing.T) {
	m := New("test-secret")

	for _, token := range []string{"", "not.a.jwt", "a.b.c"} {
		_, err := m.Parse(token)
		assert.ErrorIs(t, err, ErrInvalidToken, "token %q should be rejected", token)
	}
}
