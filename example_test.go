package human_test

import (
	"fmt"
	"github.com/anexia-it/go-human"
	"net"
	"os"
)

// SimpleChild test struct
type SimpleChild struct {
	Name      string  // no tag
	Property1 uint64  `human:"-"`          // Ignored
	Property2 float64 `human:",omitempty"` // Omitted if empty
}

// SimpleTest test struct
type SimpleTest struct {
	Var1  string      // no tag
	Var2  int         `human:"variable_2"`
	Child SimpleChild // embedded struct
}

// MapTest test struct
type MapTest struct {
	Val1      uint64
	Map       map[string]SimpleChild
	StructMap map[SimpleChild]uint8
}

// SliceTest test struct
type SliceTest struct {
	IntSlice    []int
	StructSlice []SimpleChild
}

// MapSliceTest test struct with slice of maps
type MapSliceTest struct {
	StructMapSlice []map[string]int
}

// TextMarshalerTest test struct
type TextMarshalerTest struct {
	Ip net.IP // implements encoding.TextMarshaler interface
}

// TagFailTest test struct
type TagFailTest struct {
	Test int `human:"&§/$"` // invalid tag name
}

// AnonymousFieldTest test struct
type AnonymousFieldTest struct {
	int  // ignored
	Text string
}

// UnexportedFieldTest test struct
type UnexportedFieldTest struct {
	unexported int  // ignored
	Text string
}

// Encode test with simple test struct to test encode with ignored fields and omit if field is empty
func ExampleEncoder_Encode_simpleOmitEmpty() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}
	testStruct := SimpleTest{
		Var1: "v1",
		Var2: 2,
		Child: SimpleChild{
			Name:      "theChild",
			Property1: 3, // should be ignored
			Property2: 0, // empty, should be omitted
		},
	}

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	// Output: Var1: v1
	// variable_2: 2
	// Child:
	//   Name: theChild
}

// Encode test with simple test struct to test encode with ignored fields
func ExampleEncoder_Encode_simple() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}
	testStruct := SimpleTest{
		Var1: "v1",
		Var2: 2,
		Child: SimpleChild{
			Name:      "theChild",
			Property1: 3, // should be ignored
			Property2: 4.5,
		},
	}

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	// Output: Var1: v1
	// variable_2: 2
	// Child:
	//   Name: theChild
	//   Property2: 4.5
}

// Encode test with two maps.
// Map: map string -> struct
// StructMap: map struct -> int
func ExampleEncoder_Encode_simpleMap() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	child1 := SimpleChild{
		Name:      "Person1",
		Property2: 4.5,
		Property1: 0, // should be ignored
	}

	child2 := SimpleChild{
		Name: "Person2",
	}
	stringMap := map[string]SimpleChild{
		"One": child1,
		"Two": child2,
	}
	structMap := map[SimpleChild]uint8{
		child1: 1,
		child2: 2,
	}
	testStruct := MapTest{
		Map:       stringMap,
		StructMap: structMap,
	}

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	// Output: Val1: 0
	// Map:
	//   * One: Name: Person1
	//     Property2: 4.5
	//   * Two: Name: Person2
	// StructMap:
	//   * {Person1 0 4.5}: 1
	//   * {Person2 0 0}: 2
}

// Encode test with slice of integers and slice of structs
func ExampleEncoder_Encode_simpleSlice() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	child1 := SimpleChild{
		Name:      "Person1",
		Property2: 4.5,
		Property1: 0, // should be ignored
	}

	child2 := SimpleChild{
		Name: "Person2",
	}
	structSlice := []SimpleChild{child1, child2}
	testStruct := SliceTest{
		IntSlice:    []int{1, 2, 3, 4, 5},
		StructSlice: structSlice,
	}

	// Output: IntSlice:
	//   * 1
	//   * 2
	//   * 3
	//   * 4
	//   * 5
	// StructSlice:
	//   * Name: Person1
	//     Property2: 4.5
	//   * Name: Person2

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}

// Encode test with slice of maps and each map is of type map[string]int
func ExampleEncoder_Encode_structMapSlice() {
	enc, err := human.NewEncoder(os.Stdout, human.OptionListSymbols("+", "-"))
	if err != nil {
		return
	}

	mapSliceElement1 := map[string]int{
		"one": 1,
		"two": 2,
		"tenthousandonehundredfourtytwo": 10142,
	}
	slice := []map[string]int{mapSliceElement1, mapSliceElement1}
	testStruct := MapSliceTest{
		StructMapSlice: slice,
	}

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	//Output: StructMapSlice:
	//   +
	//     - one: 1
	//     - tenthousandonehundredfourtytwo: 10142
	//     - two: 2
	//   +
	//     - one: 1
	//     - tenthousandonehundredfourtytwo: 10142
	//     - two: 2

}

// Encode test with TextMarshaler implemented by field
func ExampleEncoder_Encode_textMarshaler() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	addr := TextMarshalerTest{
		Ip: net.ParseIP("127.0.0.1"),
	}
	if err := enc.Encode(addr); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	// Output: Ip:[49 50 55 46 48 46 48 46 49]

}

// Encode test with invalid tag name
func ExampleEncoder_Encode_tagError() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	testStruct := TagFailTest{
		Test: 1,
	}

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	// Output: ERROR: 1 error occurred:
	//
	// * Invalid tag: '&§/$'
	//
}

// Encode test with anonymous field
func ExampleEncoder_Encode_anonymousFiled() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	//anonymous int field is ignored
	testStruct := AnonymousFieldTest{
		Text: "test",
	}

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	// Output: Text: test
}

// Encode test with unexported field
func ExampleEncoder_Encode_unexportedField() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	// unexported int field is ignored
	testStruct := UnexportedFieldTest{
		Text: "test",
	}

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

	// Output: Text: test
}

// Encode test iterating slice of structs
func ExampleEncoder_Encode_iteratingSlice() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	child1 := SimpleChild{
		Name:      "Person1",
		Property2: 4.5,
		Property1: 0, // should be ignored
	}

	child2 := SimpleChild{
		Name:      "Person2",
		Property2: 4.5,
		Property1: 0, // should be ignored
	}
	
	structSlice := []SimpleChild{child1, child2}
	testStruct := SliceTest{
		StructSlice: structSlice,
	}

	// Output: Name: Person1
	// Property2: 4.5
	// Name: Person2
	// Property2: 4.5

	for _, s := range testStruct.StructSlice {
		if err := enc.Encode(s); err != nil {
			fmt.Printf("ERROR; %s\n", err)
			return
		}
	}

}
