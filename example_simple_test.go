package human_test

import (
	"fmt"
	"github.com/anexia-it/go-human"
	"net"
	"os"
)

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
	StructMap map[SimpleChild]uint8
}

type SliceTest struct {
	IntSlice    []int
	StructSlice []SimpleChild
}

type MapSliceTest struct {
	StructMapSlice []map[string]int
}

type address struct {
	Ip net.IP
}

type TagFailTest struct {
	Test int `human:"&§/$"`
}

type AnonymousFieldTest struct {
	int
	Text string
}

func ExampleEncoder_Encode_SimpleOmitEmpty() {
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

	// Output: Var1: v1
	// variable_2: 2
	// Child:
	//   Name: theChild

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}

func ExampleEncoder_Encode_Simple() {
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

	// Output: Var1: v1
	// variable_2: 2
	// Child:
	//   Name: theChild
	//   Property2: 4.5

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}

func ExampleEncoder_Encode_SimpleMap() {
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

	// Output: Val1: 0
	// Map:
	//   * One: Name: Person1
	//     Property2: 4.5
	//   * Two: Name: Person2
	// StructMap:
	//   * {Person1 0 4.5}: 1
	//   * {Person2 0 0}: 2

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}

func ExampleEncoder_Encode_SimpleSlice() {
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

func ExampleEncoder_Encode_StructMapSlice() {
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

	//Output: StructMapSlice:
	//   +
	//     - one: 1
	//     - tenthousandonehundredfourtytwo: 10142
	//     - two: 2
	//   +
	//     - one: 1
	//     - tenthousandonehundredfourtytwo: 10142
	//     - two: 2
	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

}

func ExampleEncoder_Encode_TextMarshaler() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	// Output: ip: 127.0.0.1

	addr := address{
		Ip: net.ParseIP("127.0.0.1"),
	}
	if err := enc.Encode(addr); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

}

func ExampleEncoder_Encode_MapFieldError() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	testStruct := TagFailTest{
		Test: 1,
	}

	// Output: ERROR: 1 error occurred:
	//
	// * Invalid tag: '&§/$'
	//
	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}

}

func ExampleEncoder_Encode_AnonymousFiled() {
	enc, err := human.NewEncoder(os.Stdout)
	if err != nil {
		return
	}

	//anonymous int field is ignored
	testStruct := AnonymousFieldTest{
		Text: "test",
	}
	// Output: Text: test

	if err := enc.Encode(testStruct); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
}
