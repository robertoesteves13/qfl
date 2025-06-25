package qfl

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// FilterParser parses key/value data like url.Values to support (not quite) complex
// query types. It stores the keys and types that will be looked and parsed in its
// simple filtering type, for example: `x=ge!3|lt!5` will filter for any x that
// satisfies x >= x < 5
//
// Comparators:
// - [eq]: Equals
// - lt: Less than
// - gt: Greater than
// - le: Less or equal
// - ge: Greater or equal
// - like: Searches for similar string
// - is: Comma-separated list of possible options (similar to `IN ("a", "b")` in SQL)
//
// Symbols:
// | (bar): Apply more than one filter with AND
// , (comma): Separate elements in a list
// ! (mark): Indicate start of a value
// \ (backslash): escape the character in front of it (only affects symbols)
type FilterParser struct {
	keys  []string
	types []RuleType
}

func (p *FilterParser) AddInt(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, RuleTypeInt)
}

func (p *FilterParser) AddUint(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, RuleTypeUint)
}

func (p *FilterParser) AddFloat(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, RuleTypeFloat)
}

func (p *FilterParser) AddString(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, RuleTypeString)
}

func (p *FilterParser) AddTime(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, RuleTypeTime)
}

func (p FilterParser) ParseURL(u *url.URL) (*Filter, error) {
	vals := u.Query()
	fm := &Filter{}
	for i := range p.keys {
		if !vals.Has(p.keys[i]) {
			continue
		}

		tokens := p.tokenize(vals.Get(p.keys[i]))
		if len(tokens) == 1 {
			switch p.types[i] {
			case RuleTypeFloat:
				val, err := strconv.ParseFloat(tokens[0].Value, 64)
				if err != nil {
					return nil, fmt.Errorf("value `%s` is an invalid float", tokens[0].Value)
				}
				fm.AddFloat(p.keys[i], []float64{val}, ComparasionEquals)

			case RuleTypeInt:
				val, err := strconv.ParseInt(tokens[0].Value, 10, 0)
				if err != nil {
					return nil, fmt.Errorf("value `%s` is an invalid int", tokens[0].Value)
				}
				fm.AddInt(p.keys[i], []int{int(val)}, ComparasionEquals)

			case RuleTypeString:
				fm.AddString(p.keys[i], []string{tokens[0].Value}, ComparasionEquals)

			case RuleTypeTime:
				t, err := time.Parse(time.RFC3339, tokens[0].Value)
				if err != nil {
					return nil, fmt.Errorf("value `%s` is an invalid RFC3339 time", tokens[0].Value)
				}
				fm.AddTime(p.keys[i], []time.Time{t}, ComparasionEquals)

			case RuleTypeUint:
				val, err := strconv.ParseUint(tokens[0].Value, 10, 0)
				if err != nil {
					return nil, fmt.Errorf("value `%s` is an invalid uint", tokens[0].Value)
				}
				fm.AddUint(p.keys[i], []uint{uint(val)}, ComparasionEquals)

			default:
				panic(fmt.Sprintf("unexpected pkg.RuleType: %#v", p.types[i]))
			}
			continue
		}

		lastState := tokens[0].Type
		comparasion := tokens[0].ComparasionType()
		valsIdx := []int{}

		if lastState != TokenIdentifier {
			return nil, fmt.Errorf("expected comparator, found `%s`", tokens[0].Value)
		}

		if comparasion == ComparasionInvalid {
			return nil, fmt.Errorf("expected valid comparator, got `%s`", tokens[0].Value)
		}

		tokens = append(tokens, Token{Type: TokenEnd, Value: ""})
		for j := 1; j < len(tokens); j++ {
			switch tokens[j].Type {
			case TokenMark:
				if lastState != TokenIdentifier {
					return nil, fmt.Errorf("expected comparator, got `%s`", tokens[j-1].Value)
				}
			case TokenComma:
				if lastState != TokenValue {
					return nil, fmt.Errorf("expected value, got `%s`", tokens[j-1].Value)
				} else if comparasion != ComparasionIs {
					return nil, fmt.Errorf("comma is only supported on `has` comparator")
				}

			case TokenIdentifier:
				if lastState != TokenBar {
					return nil, fmt.Errorf("expected `!`, got `%s`", tokens[j].Value)
				}

				comparasion = tokens[j].ComparasionType()
			case TokenValue:
				if lastState != TokenMark && lastState != TokenComma {
					return nil, fmt.Errorf("expected `|` or comma, got `%s`", tokens[j].Value)
				}

				valsIdx = append(valsIdx, j)
			case TokenBar, TokenEnd:
				if lastState != TokenValue {
					return nil, fmt.Errorf("expected ``, got `%s`", tokens[j].Value)
				}

				fmt.Println("TEST")

				switch p.types[i] {
				case RuleTypeFloat:
					floats := make([]float64, len(valsIdx))
					for k := range valsIdx {
						val, err := strconv.ParseFloat(tokens[valsIdx[k]].Value, 64)
						if err != nil {
							return nil, fmt.Errorf("value `%s` is an invalid float", tokens[valsIdx[k]].Value)
						}

						floats[k] = val
					}

					fm.AddFloat(p.keys[i], floats, comparasion)
				case RuleTypeInt:
					ints := make([]int, len(valsIdx))
					for k := range valsIdx {
						val, err := strconv.ParseInt(tokens[valsIdx[k]].Value, 10, 0)
						if err != nil {
							return nil, fmt.Errorf("value `%s` is an invalid int", tokens[valsIdx[k]].Value)
						}

						ints[k] = int(val)
					}

					fm.AddInt(p.keys[i], ints, comparasion)
				case RuleTypeString:
					strings := make([]string, len(valsIdx))
					for k := range valsIdx {
						strings[k] = tokens[valsIdx[k]].Value
					}

					fm.AddString(p.keys[i], strings, comparasion)
				case RuleTypeTime:
					times := make([]time.Time, len(valsIdx))
					for k := range valsIdx {
						t, err := time.Parse(time.RFC3339, tokens[valsIdx[k]].Value)
						if err != nil {
							return nil, fmt.Errorf("value `%s` is an invalid RFC3339 time", tokens[valsIdx[k]].Value)
						}

						times[k] = t
					}

					fm.AddTime(p.keys[i], times, comparasion)
				case RuleTypeUint:
					uints := make([]uint, len(valsIdx))
					for k := range valsIdx {
						val, err := strconv.ParseUint(tokens[valsIdx[k]].Value, 10, 0)
						if err != nil {
							return nil, fmt.Errorf("value `%s` is an invalid uint", tokens[valsIdx[k]].Value)
						}

						uints[k] = uint(val)
					}

					fm.AddUint(p.keys[i], uints, comparasion)
				default:
					panic(fmt.Sprintf("unexpected pkg.RuleType: %#v", p.types[i]))
				}

				// Resize to 0
				valsIdx = valsIdx[:0]
			}
			lastState = tokens[j].Type
		}
	}

	return fm, nil
}

// Tokenize a string using a sliding window to work on that slice as it expands
// until it matches one of the identifier/symbols or as a value when it ends on
// either `,` or `|`, unless if either is escaped with `\`.
func (p FilterParser) tokenize(str string) (tokens []Token) {
	var slice string
	il, ih := 0, 1
	afterMark := false
	escapeNest := 0
	for ih <= len(str) {
		slice = str[il:ih]
		switch slice {
		case "!":
			afterMark = true
			tokens = append(tokens, Token{Type: TokenMark, Value: "!"})
			il = ih
		case "lt", "gt", "le", "ge", "like", "is":
			if !afterMark {
				tokens = append(tokens, Token{Type: TokenIdentifier, Value: slice})
				il = ih
			}
		default:
			last := len(slice) - 1
			lastChar := slice[last]

			if escapeNest%2 == 0 {
				switch lastChar {
				case '|':
					afterMark = false
					if len(slice) > 1 {
						slice = slice[:last]
						tokens = append(tokens, Token{Type: TokenValue, Value: slice})
					}
					tokens = append(tokens, Token{Type: TokenBar, Value: "|"})
					il = ih
				case ',':
					if len(slice) > 1 {
						slice = slice[:last]
						tokens = append(tokens, Token{Type: TokenValue, Value: slice})
					}
					tokens = append(tokens, Token{Type: TokenComma, Value: ","})

					il = ih
				}
			}

			if lastChar == '\\' {
				escapeNest += 1
			} else {
				escapeNest = 0
			}
		}

		ih += 1
	}

	if len(slice) > 0 {
		tokens = append(tokens, Token{Type: TokenValue, Value: slice})
	}

	for i := range tokens {
		tokens[i].RemoveBackslash()
	}

	return tokens
}

type Token struct {
	Type  TokenType
	Value string
}

// RemoveBackslash removes escape character `\`, except when its escaping itself.
func (t *Token) RemoveBackslash() {
	removeStack := []int{}
	for j := range t.Value {
		if t.Value[j] == '\\' {
			last := len(removeStack) - 1
			if last >= 0 && removeStack[last] == j-1 {
				removeStack[last] = j - 1
			} else {
				removeStack = append(removeStack, j)
			}
		}
	}

	for j := len(removeStack) - 1; j >= 0; j-- {
		e := removeStack[j]
		t.Value = t.Value[0:e] + t.Value[e+1:]
	}
}

func (t Token) ComparasionType() ComparasionType {
	switch t.Value {
	case "eq":
		return ComparasionEquals
	case "ge":
		return ComparasionMoreOrEqual
	case "le":
		return ComparasionLessOrEqual
	case "gt":
		return ComparasionMoreThan
	case "lt":
		return ComparasionLessThan
	case "like":
		return ComparasionLike
	case "is":
		return ComparasionIs
	}

	return ComparasionInvalid
}

type TokenType uint

const (
	TokenMark = iota + 1
	TokenBar
	TokenComma
	TokenIdentifier
	TokenValue
	TokenEnd
)

func generateSequence(start, end int) []int {
	indices := make([]int, end-start)
	for i := 0; i < end-start; i++ {
		indices[i] = start + i
	}

	return indices
}
