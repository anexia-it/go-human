package human_test

import (
	"github.com/anexia-it/go-human"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestIsNilOrEmpty(t *testing.T) {
	zeroString := ""
	nonZeroString := "test"

	// Check if nil value returns true
	require.EqualValues(t, true, human.IsNilOrEmpty(nil, reflect.ValueOf(zeroString)))

	// Check if empty string returns true
	require.EqualValues(t, true, human.IsNilOrEmpty(zeroString, reflect.ValueOf(zeroString)))
	// Check if non-empty string returns false
	require.EqualValues(t, false, human.IsNilOrEmpty(nonZeroString, reflect.ValueOf(zeroString)))

	// Check if a pointer to an empty string returns false (because the pointer is non-nil)
	require.EqualValues(t, false, human.IsNilOrEmpty(&zeroString, reflect.ValueOf(&zeroString)))
	// Check if a pointer to an empty string returns false (because the pointer is non-nil)
	require.EqualValues(t, false, human.IsNilOrEmpty(&nonZeroString, reflect.ValueOf(&zeroString)))

	// TODO: add additional test cases for types other than string
}
