package human_test

import (
	"fmt"
	"github.com/anexia-it/go-human"
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

func ExampleEncoder_Encode() {
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
