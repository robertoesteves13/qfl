package qfl_test

import (
	"fmt"

	"github.com/robertoesteves13/qfl"
)

func ExampleSQLBuilder() {
	filter := qfl.Filter{}

	filter.AddString("name", []string{"Roberto"}, qfl.ComparasionLike)
	filter.AddUint("age", []uint{23}, qfl.ComparasionMoreThan)
	filter.AddUint("salary", []uint{3000}, qfl.ComparasionLessOrEqual)
	filter.AddString("role", []string{"Programmer", "Developer"}, qfl.ComparasionEquals)

	builder := qfl.SQLBuilder{
		Filter: filter,
		Keys: map[string]string{
			"name":   "name",
			"age":    "age",
			"salary": "salary",
			"role":   "role",
		},
		PlaceholderFormat: qfl.SQLPlaceholderDollarSign,
	}

	sql, params, err := builder.Build()
	if err != nil {
		// Treat error...
	}

	fmt.Println(sql)
	fmt.Println(params)
	// Output:
	// WHERE name LIKE $1 AND age > $2 AND salary <= $3 AND role IN ($4,$5)
	// [[Roberto%] [23] [3000] [Programmer Developer]]
}
