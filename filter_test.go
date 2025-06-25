package qfl_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/robertoesteves13/qfl"

	"github.com/stretchr/testify/assert"
)

func TestParseSimpleOneValue(t *testing.T) {
	u, err := url.Parse("http://localhost:8080/api/v1/users?name=roberto")
	assert.NoError(t, err)

	parser := qfl.FilterParser{}
	parser.AddString("name")

	fm, err := parser.ParseURL(u)
	assert.NoError(t, err)
	assert.NotNil(t, fm)

	rules := fm.GetString("name")
	if assert.Equal(t, 1, len(rules)) {
		assert.Equal(t, "roberto", rules[0].Values[0])
		assert.Equal(t, qfl.ComparasionEquals, rules[0].Comparasion)
	}
}

func TestParseSimpleFourValue(t *testing.T) {
	u, err := url.Parse("http://localhost:8080/api/v1/movies?name=matrix&rate=9.2&release_date=1999-03-31T00:00:00Z&views=42")
	assert.NoError(t, err)

	parser := qfl.FilterParser{}
	parser.AddString("name")
	parser.AddFloat("rate")
	parser.AddTime("release_date")
	parser.AddUint("views")

	fm, err := parser.ParseURL(u)
	assert.NoError(t, err)
	assert.NotNil(t, fm)

	nameRule := fm.GetString("name")
	if assert.Equal(t, 1, len(nameRule)) {
		assert.Equal(t, "matrix", nameRule[0].Values[0])
		assert.Equal(t, qfl.ComparasionEquals, nameRule[0].Comparasion)
	}

	rateRule := fm.GetFloat("rate")
	if assert.Equal(t, 1, len(rateRule)) {
		assert.Equal(t, 9.2, rateRule[0].Values[0])
		assert.Equal(t, qfl.ComparasionEquals, rateRule[0].Comparasion)
	}

	releaseRule := fm.GetTime("release_date")
	if assert.Equal(t, 1, len(releaseRule)) {
		expected, _ := time.Parse(time.RFC3339, "1999-03-31T00:00:00Z")

		assert.Equal(t, expected, releaseRule[0].Values[0])
		assert.Equal(t, qfl.ComparasionEquals, releaseRule[0].Comparasion)
	}

	viewsRule := fm.GetUint("views")
	if assert.Equal(t, 1, len(viewsRule)) {
		assert.EqualValues(t, 42, viewsRule[0].Values[0])
		assert.Equal(t, qfl.ComparasionEquals, viewsRule[0].Comparasion)
	}
}

func TestParseComplexOneValue(t *testing.T) {
	u, err := url.Parse("http://localhost:8080/api/v1/employees?salary=gt!1000.0|lt!10000.0&weekHours=lt!40")
	assert.NoError(t, err)

	parser := qfl.FilterParser{}
	parser.AddFloat("salary")
	parser.AddUint("weekHours")

	fm, err := parser.ParseURL(u)
	assert.NoError(t, err)
	assert.NotNil(t, fm)

	rules := fm.GetFloat("salary")
	if assert.Equal(t, 2, len(rules)) {
		assert.Equal(t, 1000.0, rules[0].Values[0])
		assert.Equal(t, qfl.ComparasionMoreThan, rules[0].Comparasion)

		assert.Equal(t, 10000.0, rules[1].Values[0])
		assert.Equal(t, qfl.ComparasionLessThan, rules[1].Comparasion)
	}

	whRule := fm.GetUint("weekHours")
	if assert.Equal(t, 1, len(whRule)) {
		assert.EqualValues(t, 40, whRule[0].Values[0])
		assert.Equal(t, qfl.ComparasionLessThan, whRule[0].Comparasion)
	}
}

func TestParserIs(t *testing.T) {
	u, err := url.Parse("http://localhost:8080/api/v1/employees?role=is!Programmer,Tester,DBA")
	assert.NoError(t, err)

	parser := qfl.FilterParser{}
	parser.AddString("role")

	fm, err := parser.ParseURL(u)
	assert.NoError(t, err)
	assert.NotNil(t, fm)

	rules := fm.GetString("role")
	if assert.Equal(t, 1, len(rules)) {
		assert.Equal(t, qfl.ComparasionIs, rules[0].Comparasion)
		assert.Equal(t, "Programmer", rules[0].Values[0])
		assert.Equal(t, "Tester", rules[0].Values[1])
		assert.Equal(t, "DBA", rules[0].Values[2])

	}
}
