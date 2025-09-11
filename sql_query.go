package qfl

import (
	"fmt"
	"strconv"
	"strings"
)

type SQLPlaceholderFormat uint8

const (
	SQLPlaceholderQuestionMark SQLPlaceholderFormat = 0
	SQLPlaceholderDollarSign   SQLPlaceholderFormat = 1
)

// SQLBuilder is a generic WHERE-condition builder capable to convert filter
// rules automatically.
//
// All the functions build the string on the builder in place, and it doesn't
// check nor support if you call them in the wrong order. Use with caution
type SQLBuilder struct {
	// string builder for the SQL query
	Builder strings.Builder

	Filter            Filter
	Keys              map[string]string
	PlaceholderFormat SQLPlaceholderFormat
}

func (sq *SQLBuilder) Select(table string, columns ...string) {
	sq.Builder.WriteString("SELECT ")

	for i := range columns {
		sq.Builder.WriteString(columns[i])

		if i == len(columns)-1 {
			sq.Builder.WriteRune(' ')
		} else {
			sq.Builder.WriteString(", ")
		}
	}

	sq.Builder.WriteString("FROM ")
	sq.Builder.WriteString(table)
	sq.Builder.WriteRune('\n')
}

func (sq *SQLBuilder) Where() ([]any, error) {
	if sq.Keys == nil {
		return nil, fmt.Errorf("field `Keys` is empty")
	}

	parameters := []any{}
	offset := uint(0)

	sq.Builder.WriteString("WHERE ")
	for i := range sq.Filter.keys {
		key := sq.Filter.keys[i]
		if column, ok := sq.Keys[key.key]; ok {
			for i := range key.rules {
				var (
					params []any
				)

				switch key.Type {
				case ruleTypeInt:
					params = extractConditions(column, key.rules[i], sq.Filter.intVals, offset, sq.PlaceholderFormat, &sq.Builder)
				case ruleTypeUint:
					params = extractConditions(column, key.rules[i], sq.Filter.uintVals, offset, sq.PlaceholderFormat, &sq.Builder)
				case ruleTypeFloat:
					params = extractConditions(column, key.rules[i], sq.Filter.floatVals, offset, sq.PlaceholderFormat, &sq.Builder)
				case ruleTypeString:
					params = extractConditions(column, key.rules[i], sq.Filter.stringVals, offset, sq.PlaceholderFormat, &sq.Builder)
				case ruleTypeTime:
					params = extractConditions(column, key.rules[i], sq.Filter.timeVals, offset, sq.PlaceholderFormat, &sq.Builder)
				}

				parameters = append(parameters, params...)
				offset += uint(len(params))

				if i != len(key.rules)-1 {
					sq.Builder.WriteString(" AND ")
				}
			}
		}

		if i != len(sq.Filter.keys)-1 {
			sq.Builder.WriteString(" AND ")
		}
	}

	sq.Builder.WriteRune('\n')
	return parameters, nil
}

func (sq *SQLBuilder) Join(table, condition string) {
	sq.Builder.WriteString("JOIN ")
	sq.Builder.WriteString(table)
	sq.Builder.WriteString(" ON ")
	sq.Builder.WriteString(condition)
	sq.Builder.WriteRune('\n')
}

// Page paginates the query by the limit number. Pages are zero-indexed
func (sq *SQLBuilder) Page(limit, page uint64) {
	sq.Builder.WriteString("LIMIT ")
	sq.Builder.Write(strconv.AppendUint(nil, limit, 10))
	sq.Builder.WriteString(" OFFSET ")
	sq.Builder.Write(strconv.AppendUint(nil, page*limit, 10))
	sq.Builder.WriteRune('\n')
}

func (sq *SQLBuilder) Order(order string, columns ...string) {
	sq.Builder.WriteString("ORDER BY ")

	for i := range columns {
		sq.Builder.WriteString(columns[i])

		if i == len(columns)-1 {
			sq.Builder.WriteRune(' ')
		} else {
			sq.Builder.WriteString(", ")
		}
	}

	sq.Builder.WriteString(order)
	sq.Builder.WriteRune('\n')
}

func (sq *SQLBuilder) With(name string, statement string) {
	sq.Builder.WriteString("WITH ")
	sq.Builder.WriteString(name)
	sq.Builder.WriteString(" AS ( ")
	sq.Builder.WriteString(statement)
	sq.Builder.WriteRune(')')
	sq.Builder.WriteRune('\n')
}

func extractConditions[T Primitive](column string, rule filterRule, values []T, offset uint, format SQLPlaceholderFormat, builder *strings.Builder) (params []any) {
	params = make([]any, len(rule.indices))
	for i := range rule.indices {
		params[i] = values[rule.indices[i]]
	}

	skipPlaceholder := false

	builder.WriteString(column)
	switch rule.Comparasion {
	case ComparasionEquals:
		if len(params) > 1 {
			builder.WriteString(" IN ")
			stringifyListParams(params, offset, format, builder)
			skipPlaceholder = true
		} else {
			builder.WriteString(" = ")
		}

	case ComparasionLessThan:
		builder.WriteString(" < ")
	case ComparasionMoreThan:
		builder.WriteString(" > ")
	case ComparasionLessOrEqual:
		builder.WriteString(" <= ")
	case ComparasionMoreOrEqual:
		builder.WriteString(" >= ")
	case ComparasionLike:
		builder.WriteString(" LIKE ")
		for i := range params {
			params[i] = fmt.Sprint(params[rule.indices[i]])
		}
	}

	if !skipPlaceholder && format == SQLPlaceholderDollarSign {
		builder.WriteRune('$')
		builder.Write(strconv.AppendInt(nil, int64(1+offset), 10))
	} else if !skipPlaceholder {
		builder.WriteRune('?')
	}

	return
}

func stringifyListParams(params []any, offset uint, format SQLPlaceholderFormat, builder *strings.Builder) {
	builder.WriteRune('(')

	for i := range params {
		if format == SQLPlaceholderDollarSign {
			builder.WriteRune('$')
			builder.Write(strconv.AppendInt(nil, int64(1+offset)+int64(i), 10))
		} else {
			builder.WriteRune('?')
		}

		if i != len(params)-1 {
			builder.WriteRune(',')
		}
	}

	builder.WriteRune(')')
}
