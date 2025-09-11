package qfl_test

import (
	"fmt"

	"github.com/robertoesteves13/qfl"
)

func ExampleSQLBuilder() {
	filter := qfl.Filter{}

	filter.AddString("name", []string{"Roberto%"}, qfl.ComparasionLike)
	filter.AddUint("age", []uint{23}, qfl.ComparasionMoreThan)
	filter.AddUint("age", []uint{60}, qfl.ComparasionLessOrEqual)
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

	builder.Select("employer", "id", "name", "age", "salary", "role", "employed_since")
	params, err := builder.Where()
	builder.Page(20, 5)

	sql := builder.Builder.String()
	if err != nil {
		// Treat error...
	}

	fmt.Println(sql)
	fmt.Println(params)
	// Output:
	// SELECT id, name, age, salary, role, employed_since FROM employer
	// WHERE name LIKE $1 AND age > $2 AND age <= $3 AND salary <= $4 AND role IN ($5,$6)
	// LIMIT 20 OFFSET 100
	//
	// [Roberto% 23 60 3000 Programmer Developer]
}
