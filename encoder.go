package human

import (
	"encoding"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/speijnik/go-errortree"
	"io"
	"reflect"
	"strings"
	"unicode"
)

// Encoder writes human readable text to an output stream.
type Encoder struct {
	stream      *FlushableBuffer
	tagName     string
	indent      uint
	listSymbols []string
}

// Encode writes the human encoding of v to the stream.
func (e *Encoder) Encode(v interface{}) error {
	value := reflect.ValueOf(v)
	if err := e.encodeStruct(value, 0); err != nil {
		e.stream.Reset()
		return err
	}
	_, err := e.stream.Flush()
	if err != nil {
		e.stream.Reset()
	}
	return err
}

func (e *Encoder) encodeStruct(v reflect.Value, indentLevel uint) (err error) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldDefinition := t.Field(i)
		if fieldDefinition.Anonymous {
			// skip anonymous field
			continue
		}
		fieldValue := v.Field(i)

		fieldName := fieldDefinition.Name

		if !unicode.IsUpper([]rune(fieldName)[0]) {
			// Ignore private fields
			continue
		}

		fieldName, omitEmpty, tagErr := parseTagFromStructField(fieldDefinition, e.tagName)
		if tagErr != nil {
			// Parsing the tag failed, ignore the field and carry on
			err = multierror.Append(err, tagErr)
			continue
		}

		fieldInterface := fieldValue.Interface()
		if fieldInterface == nil || fieldName == "-" || (omitEmpty && IsNilOrEmpty(fieldInterface, fieldValue)) {
			// Skip field if:
			// - field is a nil-value
			// - field name specifies that the field shall be omitted
			// - omitEmpty is set and the field is nil or empty
			continue
		}

		// Getting this far means we are handling a non-empty field
		fmt.Fprint(e.stream, strings.Repeat(" ", int(e.indent*indentLevel))+fieldName+":")
		if fieldEncodeErr := e.encodeValue(fieldInterface, fieldValue, indentLevel); fieldEncodeErr != nil {
			err = errortree.Add(err, fieldName, fieldEncodeErr)
		}
	}

	return
}

func (e *Encoder) encodeSlice(v reflect.Value, indentLevel uint) error {

	listSymbol := e.listSymbols[int(indentLevel)%len(e.listSymbols)]

	for i := 0; i < v.Len(); i++ {
		valueV := v.Index(i)
		valueI := valueV.Interface()
		fmt.Fprint(e.stream, strings.Repeat(" ", int(e.indent*indentLevel))+listSymbol+" ")
		if err := e.encodeValue(valueI, valueV, indentLevel); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeMap(v reflect.Value, indentLevel uint) error {

	listSymbol := e.listSymbols[int(indentLevel)%len(e.listSymbols)]

	keys := v.MapKeys()
	for i := 0; i < len(keys); i++ {
		keyV := keys[i]
		keyI := keyV.Interface()
		valueV := v.MapIndex(keyV)
		fmt.Fprint(e.stream, strings.Repeat(" ", int(e.indent*indentLevel))+listSymbol+" "+fmt.Sprint(keyI)+": ")
		if err := e.encodeValue(valueV.Interface(), valueV, indentLevel); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeValue(i interface{}, v reflect.Value, indentLevel uint) (err error) {
	// Check if the passed interface implements encoding.TextMarshaler, in which case we use the marshaler
	// for generating the value
	if marshaler, ok := i.(encoding.TextMarshaler); ok {
		text, marshalErr := marshaler.MarshalText()
		if marshalErr != nil {
			err = marshalErr
		}
		// As MarshalText is expected to return a textual representation, print it to our stream
		fmt.Fprintln(e.stream, text)
		return
	}

	// At this point it is safe to get rid of a possible pointer...
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	} else if v.Kind() == reflect.Ptr {
		// No-op for nil-pointers
		return
	}

	// Per-type handling
	switch v.Kind() {
	case reflect.Struct:
		// Handle struct
		fmt.Fprintln(e.stream, "")
		err = e.encodeStruct(v, indentLevel+1)
	case reflect.Slice, reflect.Array:
		// Handle slice
		fmt.Fprintln(e.stream, "")
		err = e.encodeSlice(v, indentLevel+1)
	case reflect.Map:
		// Handle map
		fmt.Fprintln(e.stream, "")
		err = e.encodeMap(v, indentLevel+1)

	default:
		// All other types are mapped as-is
		// missuse Fprint's sepereration spaces to introduce a space in front of the value
		fmt.Fprintln(e.stream, "", i)
	}

	return
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer, opts ...Option) (encoder *Encoder, err error) {
	encoder = &Encoder{
		stream: NewFlushableBuffer(w),
	}

	// apply options
	for _, opt := range append(defaultOptions, opts...) {
		if optErr := opt(encoder); optErr != nil {
			err = multierror.Append(err, optErr)
		}
	}
	// check if any option returned an error and set encoder nil
	if err != nil {
		encoder = nil
	}
	return
}
