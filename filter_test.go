package qfl_test

import (
	"fmt"

	"github.com/robertoesteves13/qfl"
)

func ExampleFilter() {
	f := qfl.Filter{}

	// Set rules inside filter
	f.AddInt("age", []int{22}, qfl.ComparasionMoreThan)
	f.AddString("name", []string{"John"}, qfl.ComparasionLike)

	// Get filter rules
	age := f.GetInt("age")
	name := f.GetString("name")

	fmt.Printf("age %s %d\n", age[0].Comparasion, age[0].Values[0])
	fmt.Printf("name %s %s\n", name[0].Comparasion, name[0].Values[0])
	// Output:
	// age MoreThan 22
	// name Like John
}
