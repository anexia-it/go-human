package human

import (
	"bytes"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewEncoder(t *testing.T) {
	t.Run("DefaultOptions", func(t *testing.T) {
		enc, err := NewEncoder(nil)
		require.NoError(t, err)
		require.NotNil(t, enc)
		require.EqualValues(t, DefaultTagName, enc.tagName)
		require.EqualValues(t, []string{DefaultListSymbol}, enc.listSymbols)
		require.EqualValues(t, DefaultIndent, enc.indent)
	})

	t.Run("CustomOptionsNoError", func(t *testing.T) {
		enc, err := NewEncoder(nil,
			OptionListSymbols("+", "-"),
			OptionTagName("test"),
			OptionIndent(4))
		require.NoError(t, err)
		require.NotNil(t, enc)
		require.EqualValues(t, "test", enc.tagName)
		require.EqualValues(t, []string{"+", "-"}, enc.listSymbols)
		require.EqualValues(t, 4, enc.indent)
	})

	t.Run("CustomOptionsError", func(t *testing.T) {
		enc, err := NewEncoder(nil,
			OptionListSymbols(),
			OptionTagName(""))
		require.Error(t, err)
		require.Nil(t, enc)
		multiErr, ok := err.(*multierror.Error)
		require.True(t, ok, "Must be a *multierror.Error")
		require.Len(t, multiErr.Errors, 2)
		require.EqualError(t, multiErr.Errors[0], ErrListSymbolsEmpty.Error())
		require.EqualError(t, multiErr.Errors[1], ErrInvalidTagName.Error())
	})

}
