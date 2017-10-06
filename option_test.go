package human

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOptionTagName(t *testing.T) {

	enc := &Encoder{}

	opt := OptionTagName("test")

	require.NoError(t, opt(enc))
	require.EqualValues(t, "test", enc.tagName)
}

func TestOptionListSymbols(t *testing.T) {

	enc := &Encoder{}

	opt := OptionListSymbols("*", "+", "-")

	require.NoError(t, opt(enc))
	require.EqualValues(t, []string{"*", "+", "-"}, enc.listSymbols)
}

func TestOptionIndent(t *testing.T) {

	enc := &Encoder{}

	opt := OptionIndent(4)

	require.NoError(t, opt(enc))
	require.EqualValues(t, 4, enc.indent)
}
