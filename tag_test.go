package human

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsInvalidTag(t *testing.T) {
	err := newErrorInvalidTag("testing error")
	tag, isInvalid := IsInvalidTag(err)
	require.NotNil(t, tag)
	require.True(t, isInvalid)
}

func TestInvalidTagError(t *testing.T) {
	err := newErrorInvalidTag("testing error")
	require.EqualValues(t, "Invalid tag: 'testing error'", err.Error())
}

func TestInvalidTagTag(t *testing.T) {
	err := newErrorInvalidTag("testing error")
	tag, isInvalid := IsInvalidTag(err)
	require.True(t, isInvalid)
	require.EqualValues(t, "testing error", tag.Tag())
}

func TestParseTag(t *testing.T) {
	name, omitEmpty, err := ParseTag("test")
	require.NoError(t, err)
	require.False(t, omitEmpty)
	require.EqualValues(t, "test", name)
}

func TestParseTagError(t *testing.T) {
	_, _, err := ParseTag("&")
	require.Error(t, err)
	tag, isInvalid := IsInvalidTag(err)
	require.True(t, isInvalid)
	require.EqualValues(t, "&", tag.Tag())
}
