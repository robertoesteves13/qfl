package qfl

import (
	"fmt"
	"strings"
)

type SQLPlaceholderFormat uint8

const (
	SQLPlaceholderQuestionMark SQLPlaceholderFormat = 0
	SQLPlaceholderDollarSign   SQLPlaceholderFormat = 1
)

// SQLBuilder is a generic WHERE-condition builder capable to convert filter
// rules automatically
type SQLBuilder struct {
	Filter            Filter
	Keys              map[string]string
	PlaceholderFormat SQLPlaceholderFormat
}

func (sq *SQLBuilder) Build() (string, []any, error) {
	if sq.Keys == nil {
		return "", nil, fmt.Errorf("field `Keys` is empty")
	}

	conditions := []string{}
	parameters := []any{}
	offset := uint(0)

	for i := range sq.Filter.keys {
		key := sq.Filter.keys[i]
		if column, ok := sq.Keys[key.key]; ok {
			for i := range key.rules {
				var (
					expr   string
					params []any
				)

				switch key.Type {
				case ruleTypeInt:
					expr, params = extractConditions(column, key.rules[i], sq.Filter.intVals, offset, sq.PlaceholderFormat)
				case ruleTypeUint:
					expr, params = extractConditions(column, key.rules[i], sq.Filter.uintVals, offset, sq.PlaceholderFormat)
				case ruleTypeFloat:
					expr, params = extractConditions(column, key.rules[i], sq.Filter.floatVals, offset, sq.PlaceholderFormat)
				case ruleTypeString:
					expr, params = extractConditions(column, key.rules[i], sq.Filter.stringVals, offset, sq.PlaceholderFormat)
				case ruleTypeTime:
					expr, params = extractConditions(column, key.rules[i], sq.Filter.timeVals, offset, sq.PlaceholderFormat)
				}

				conditions = append(conditions, expr)
				parameters = append(parameters, params)
				offset += uint(len(params))
			}
		}
	}

	return "WHERE " + strings.Join(conditions, " AND "), parameters, nil
}

func extractConditions[T Primitive](column string, rule filterRule, values []T, offset uint, format SQLPlaceholderFormat) (expression string, params []any) {
	params = make([]any, len(rule.indices))
	for i := range rule.indices {
		params[i] = values[rule.indices[i]]
	}

	skipPlaceholder := true

	switch rule.Comparasion {
	case ComparasionEquals:
		if len(params) > 1 {
			expression = column + " IN " + stringifyListParams(params, offset, format)
		} else {
			expression = fmt.Sprintf("%s = $%d", column, 1+offset)
		}

		skipPlaceholder = false
	case ComparasionLessThan:
		expression = fmt.Sprintf("%s < ", column)
	case ComparasionMoreThan:
		expression = fmt.Sprintf("%s > ", column)
	case ComparasionLessOrEqual:
		expression = fmt.Sprintf("%s <= ", column)
	case ComparasionMoreOrEqual:
		expression = fmt.Sprintf("%s >= ", column)
	case ComparasionLike:
		expression = fmt.Sprintf("%s LIKE ", column)
		for i := range params {
			params[i] = fmt.Sprint(params[rule.indices[i]]) + "%"
		}
	}

	if skipPlaceholder && format == SQLPlaceholderDollarSign {
		expression += fmt.Sprintf("$%d", 1+offset)
	} else if skipPlaceholder {
		expression += "?"
	}

	return
}

// PERF: Benchmark this
func stringifyListParams(params []any, offset uint, format SQLPlaceholderFormat) string {
	strs := make([]string, len(params))
	for i := range strs {
		if format == SQLPlaceholderDollarSign {
			strs[i] = fmt.Sprintf("$%d", uint(i+1)+offset)
		} else {
			strs[i] = "?"
		}
	}

	return "(" + strings.Join(strs, ",") + ")"
}
