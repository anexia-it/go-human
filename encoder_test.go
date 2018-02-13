package human

import (
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
	"bytes"
	"github.com/stretchr/testify/assert"
	"time"
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

func TestEncoder_Encode(t *testing.T) {
	outputBuffer := bytes.NewBufferString("")
	enc, err := NewEncoder(outputBuffer)
	require.NoError(t, err)
	require.NotNil(t, enc)

	t.Run("Map", func(t *testing.T) {
		outputBuffer.Reset()

		m := map[string]interface{}{
			"test0": "0",
			"test1": 1,
			"test2": 2.9,
			"test3": time.Second * 3,
		}

		expectedOutput := "\n* test0: 0\n* test1: 1\n* test2: 2.9\n* test3: 3s\n"

		assert.NoError(t, enc.Encode(m))
		assert.EqualValues(t, expectedOutput, outputBuffer.String())
	})

	t.Run("Slice", func(t *testing.T) {
		outputBuffer.Reset()

		s := []interface{}{
			"test0", 1, 2.9, time.Second * 3,
		}

		expectedOutput := "\n* test0\n* 1\n* 2.9\n* 3s\n"

		assert.NoError(t, enc.Encode(s))
		assert.EqualValues(t, expectedOutput, outputBuffer.String())
	})
}
