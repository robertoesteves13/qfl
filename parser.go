package qfl

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Parser parses the QFL language for the keys you specify
type Parser struct {
	TimeFormat string // defaults to RCF3339 if empty
	keys       []string
	types      []ruleType
}

func (p *Parser) AddInt(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, ruleTypeInt)
}

func (p *Parser) AddUint(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, ruleTypeUint)
}

func (p *Parser) AddFloat(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, ruleTypeFloat)
}

func (p *Parser) AddString(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, ruleTypeString)
}

func (p *Parser) AddTime(key string) {
	p.keys = append(p.keys, key)
	p.types = append(p.types, ruleTypeTime)
}

// ParseURL reads query variables and returns the filter containing all rules
// for them. Note that it only looks for the first definition of the variable.
func (p Parser) ParseURL(u *url.URL) (*Filter, error) {
	vals := u.Query()
	kv := make(map[string]string)
	for k, v := range vals {
		kv[k] = v[0]
	}

	return p.Parse(kv)
}

func (p Parser) Parse(kv map[string]string) (*Filter, error) {
	// Ensure that time format is set before running the parser
	if p.TimeFormat == "" {
		p.TimeFormat = time.RFC3339
	}

	fm := &Filter{}
	for i := range p.keys {
		if _, ok := kv[p.keys[i]]; !ok {
			continue
		}

		tokens := p.tokenize(kv[p.keys[i]])
		if len(tokens) == 1 {
			switch p.types[i] {
			case ruleTypeFloat:
				val, err := strconv.ParseFloat(tokens[0].Value, 64)
				if err != nil {
					return nil, fmt.Errorf("value `%s` is an invalid float", tokens[0].Value)
				}
				fm.AddFloat(p.keys[i], []float64{val}, ComparasionEquals)

			case ruleTypeInt:
				val, err := strconv.ParseInt(tokens[0].Value, 10, 0)
				if err != nil {
					return nil, fmt.Errorf("value `%s` is an invalid int", tokens[0].Value)
				}
				fm.AddInt(p.keys[i], []int{int(val)}, ComparasionEquals)

			case ruleTypeString:
				fm.AddString(p.keys[i], []string{tokens[0].Value}, ComparasionEquals)

			case ruleTypeTime:
				t, err := time.Parse(p.TimeFormat, tokens[0].Value)
				if err != nil {
					return nil, fmt.Errorf("value `%s` is not a time formatted as `%s`", tokens[0].Value, p.TimeFormat)
				}
				fm.AddTime(p.keys[i], []time.Time{t}, ComparasionEquals)

			case ruleTypeUint:
				val, err := strconv.ParseUint(tokens[0].Value, 10, 0)
				if err != nil {
					return nil, fmt.Errorf("value `%s` is an invalid uint", tokens[0].Value)
				}
				fm.AddUint(p.keys[i], []uint{uint(val)}, ComparasionEquals)

			default:
				return nil, fmt.Errorf("unexpected pkg.RuleType: %#v", p.types[i])
			}
			continue
		}

		lastState := tokens[0].Type
		comparasion := tokens[0].comparasionType()
		valsIdx := []int{}

		if lastState != tokenIdentifier {
			return nil, fmt.Errorf("expected comparator, found `%s`", tokens[0].Value)
		}

		if comparasion == ComparasionInvalid {
			return nil, fmt.Errorf("expected valid comparator, got `%s`", tokens[0].Value)
		}

		tokens = append(tokens, token{Type: tokenEnd, Value: ""})
		for j := 1; j < len(tokens); j++ {
			switch tokens[j].Type {
			case tokenMark:
				if lastState != tokenIdentifier {
					return nil, fmt.Errorf("expected comparator, got `%s`", tokens[j-1].Value)
				}
			case tokenComma:
				if lastState != tokenValue {
					return nil, fmt.Errorf("expected value, got `%s`", tokens[j-1].Value)
				} else if comparasion != ComparasionEquals {
					return nil, fmt.Errorf("comma is only supported on `eq` comparator")
				}

			case tokenIdentifier:
				if lastState != tokenBar {
					return nil, fmt.Errorf("expected `!`, got `%s`", tokens[j].Value)
				}

				comparasion = tokens[j].comparasionType()
			case tokenValue:
				if lastState != tokenMark && lastState != tokenComma {
					return nil, fmt.Errorf("expected `|` or comma, got `%s`", tokens[j].Value)
				}

				valsIdx = append(valsIdx, j)
			case tokenBar, tokenEnd:
				if lastState != tokenValue {
					return nil, fmt.Errorf("expected ``, got `%s`", tokens[j].Value)
				}

				switch p.types[i] {
				case ruleTypeFloat:
					floats := make([]float64, len(valsIdx))
					for k := range valsIdx {
						val, err := strconv.ParseFloat(tokens[valsIdx[k]].Value, 64)
						if err != nil {
							return nil, fmt.Errorf("value `%s` is an invalid float", tokens[valsIdx[k]].Value)
						}

						floats[k] = val
					}

					fm.AddFloat(p.keys[i], floats, comparasion)
				case ruleTypeInt:
					ints := make([]int, len(valsIdx))
					for k := range valsIdx {
						val, err := strconv.ParseInt(tokens[valsIdx[k]].Value, 10, 0)
						if err != nil {
							return nil, fmt.Errorf("value `%s` is an invalid int", tokens[valsIdx[k]].Value)
						}

						ints[k] = int(val)
					}

					fm.AddInt(p.keys[i], ints, comparasion)
				case ruleTypeString:
					strings := make([]string, len(valsIdx))
					for k := range valsIdx {
						strings[k] = tokens[valsIdx[k]].Value
					}

					fm.AddString(p.keys[i], strings, comparasion)
				case ruleTypeTime:
					times := make([]time.Time, len(valsIdx))
					for k := range valsIdx {
						t, err := time.Parse(p.TimeFormat, tokens[valsIdx[k]].Value)
						if err != nil {
							return nil, fmt.Errorf("value `%s` is not a time formatted as `%s`", tokens[valsIdx[k]].Value, p.TimeFormat)
						}

						times[k] = t
					}

					fm.AddTime(p.keys[i], times, comparasion)
				case ruleTypeUint:
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
					return nil, fmt.Errorf("unexpected pkg.RuleType: %#v", p.types[i])
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
func (p Parser) tokenize(str string) (tokens []token) {
	var slice string
	il, ih := 0, 1
	afterMark := false
	escapeNest := 0
	for ih <= len(str) {
		slice = str[il:ih]
		switch slice {
		case "!":
			afterMark = true
			tokens = append(tokens, token{Type: tokenMark, Value: "!"})
			il = ih
		case "lt", "gt", "le", "ge", "lk", "eq":
			if !afterMark {
				tokens = append(tokens, token{Type: tokenIdentifier, Value: slice})
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
						tokens = append(tokens, token{Type: tokenValue, Value: slice})
					}
					tokens = append(tokens, token{Type: tokenBar, Value: "|"})
					il = ih
				case ',':
					if len(slice) > 1 {
						slice = slice[:last]
						tokens = append(tokens, token{Type: tokenValue, Value: slice})
					}
					tokens = append(tokens, token{Type: tokenComma, Value: ","})

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
		tokens = append(tokens, token{Type: tokenValue, Value: slice})
	}

	// Do another pass to remove the escaped symbol
	for i := range tokens {
		tokens[i].removeBackslash()
	}

	return tokens
}

type token struct {
	Type  tokenType
	Value string
}

// RemoveBackslash removes escape character `\`, except when its escaping itself.
func (t *token) removeBackslash() {
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

func (t token) comparasionType() ComparasionType {
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
	case "lk":
		return ComparasionLike
	}

	return ComparasionInvalid
}

type tokenType uint

const (
	tokenMark = iota + 1
	tokenBar
	tokenComma
	tokenIdentifier
	tokenValue
	tokenEnd
)
