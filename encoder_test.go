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

type SimpleChild struct {
	Name      string  // no tag
	Property1 uint64  `human:"-"`          // Ignored
	Property2 float64 `human:",omitempty"` // Omitted if empty
}
type SimpleTest struct {
	Var1  string //no tag
	Var2  int    `human:"variable_2"`
	Child SimpleChild
}

type MapTest struct {
	Val1      uint64
	Map       map[string]SimpleChild
	structMap map[SimpleChild]uint8
}

type SliceTest struct {
	intSlice    []int
	structSlice []SimpleChild
}

func TestEncoder_Encode(t *testing.T) {
	t.Run("SimpleWithOmitEmpty", func(t *testing.T) {
		buf := bytes.NewBufferString("")
		enc, err := NewEncoder(buf)
		require.NoError(t, err)
		require.NotNil(t, enc)

		testStruct := SimpleTest{
			Var1: "v1",
			Var2: 2,
			Child: SimpleChild{
				Name:      "theChild",
				Property1: 3, // should be ignored
				Property2: 0, // empty, should be omitted
			},
		}

		expected := `Var1: v1
variable_2: 2
Child:
  Name: theChild
`
		require.NoError(t, enc.Encode(testStruct))
		require.EqualValues(t, expected, buf.String())

	})

	t.Run("Simple", func(t *testing.T) {
		buf := bytes.NewBufferString("")
		enc, err := NewEncoder(buf)
		require.NoError(t, err)
		require.NotNil(t, enc)

		testStruct := SimpleTest{
			Var1: "v1",
			Var2: 2,
			Child: SimpleChild{
				Name:      "theChild",
				Property1: 3,   // should be ignored
				Property2: 4.5, // not empty, should be encoded
			},
		}

		expected := `Var1: v1
variable_2: 2
Child:
  Name: theChild
  Property2: 4.5
`
		require.NoError(t, enc.Encode(testStruct))
		require.EqualValues(t, expected, buf.String())

	})
}
