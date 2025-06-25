package qfl

import (
	"time"
)

// Primitive is a generic interface that indicates the types the filter can store.
type Primitive interface {
	int | uint | float64 | string | time.Time
}

// ComparasionType indicates the comparasion it should make for the values.
type ComparasionType uint8

const (
	ComparasionInvalid ComparasionType = iota
	ComparasionEquals
	ComparasionLessThan
	ComparasionMoreThan
	ComparasionLessOrEqual
	ComparasionMoreOrEqual
	ComparasionLike
)

func (c ComparasionType) String() string {
	switch c {
	case ComparasionEquals:
		return "Equals"
	case ComparasionLessOrEqual:
		return "LessOrEqual"
	case ComparasionLessThan:
		return "LessThan"
	case ComparasionLike:
		return "Like"
	case ComparasionMoreOrEqual:
		return "MoreOrEqual"
	case ComparasionMoreThan:
		return "MoreThan"
	default:
		return "Invalid"
	}
}

// FilterRule represents
// Note that `ComparasionEquals` is the only one that can have more than one value.
type FilterRule[T Primitive] struct {
	Comparasion ComparasionType
	Values      []T
}

// Filter is a specialized data structure that stores rules for a given key. It
// only supports some primitive data types. All get and set functions should be
// the exactly same except for the type it's manipulating, this is on purporse
// to ensure type safety.
type Filter struct {
	keys []filterKey

	intVals    []int
	uintVals   []uint
	floatVals  []float64
	stringVals []string
	timeVals   []time.Time
}

func (f *Filter) GetInt(key string) []FilterRule[int] {
	for i := range f.keys {
		if f.keys[i].key == key && f.keys[i].Type == ruleTypeInt {
			return getGeneric(f.keys[i], f.intVals)
		}
	}

	return nil
}

func (f *Filter) GetUint(key string) []FilterRule[uint] {
	for i := range f.keys {
		if f.keys[i].key == key && f.keys[i].Type == ruleTypeUint {
			return getGeneric(f.keys[i], f.uintVals)
		}
	}

	return nil
}

func (f *Filter) GetFloat(key string) []FilterRule[float64] {
	for i := range f.keys {
		if f.keys[i].key == key && f.keys[i].Type == ruleTypeFloat {
			return getGeneric(f.keys[i], f.floatVals)
		}
	}

	return nil
}

func (f *Filter) GetString(key string) []FilterRule[string] {
	for i := range f.keys {
		if f.keys[i].key == key && f.keys[i].Type == ruleTypeString {
			return getGeneric(f.keys[i], f.stringVals)
		}
	}

	return nil
}

func (f *Filter) GetTime(key string) []FilterRule[time.Time] {
	for i := range f.keys {
		if f.keys[i].key == key && f.keys[i].Type == ruleTypeTime {
			return getGeneric(f.keys[i], f.timeVals)
		}
	}

	return nil
}

func getGeneric[T Primitive](key filterKey, vals []T) []FilterRule[T] {
	rules := key.rules
	rulesReturn := make([]FilterRule[T], len(rules))

	for j := range rules {
		indices := rules[j].indices
		values := make([]T, len(indices))
		for k := range indices {
			values[k] = vals[indices[k]]
		}

		rulesReturn[j] = FilterRule[T]{
			Comparasion: rules[j].Comparasion,
			Values:      values,
		}
	}

	return rulesReturn

}

func (f *Filter) AddInt(key string, values []int, comparasion ComparasionType) {
	start := len(f.intVals)
	f.intVals = append(f.intVals, values...)
	end := len(f.intVals)

	indices := generateSequence(start, end)
	f.appendRule(key, indices, comparasion, ruleTypeInt)
}

func (f *Filter) AddUint(key string, values []uint, comparasion ComparasionType) {
	start := len(f.uintVals)
	f.uintVals = append(f.uintVals, values...)
	end := len(f.uintVals)

	indices := generateSequence(start, end)
	f.appendRule(key, indices, comparasion, ruleTypeUint)
}

func (f *Filter) AddFloat(key string, values []float64, comparasion ComparasionType) {
	start := len(f.floatVals)
	f.floatVals = append(f.floatVals, values...)
	end := len(f.floatVals)

	indices := generateSequence(start, end)
	f.appendRule(key, indices, comparasion, ruleTypeFloat)
}

func (f *Filter) AddString(key string, values []string, comparasion ComparasionType) {
	start := len(f.stringVals)
	f.stringVals = append(f.stringVals, values...)
	end := len(f.stringVals)

	indices := generateSequence(start, end)
	f.appendRule(key, indices, comparasion, ruleTypeString)
}

func (f *Filter) AddTime(key string, values []time.Time, comparasion ComparasionType) {
	start := len(f.timeVals)
	f.timeVals = append(f.timeVals, values...)
	end := len(f.timeVals)

	indices := generateSequence(start, end)
	f.appendRule(key, indices, comparasion, ruleTypeTime)
}

func (f *Filter) appendRule(key string, indices []int, comparasion ComparasionType, ruleType ruleType) {
	rule := filterRule{
		Comparasion: comparasion,
		indices:     indices,
	}

	for i := range f.keys {
		if f.keys[i].key == key {
			f.keys[i].rules = append(f.keys[i].rules, rule)
			return
		}
	}

	k := filterKey{
		key:   key,
		Type:  ruleType,
		rules: []filterRule{rule},
	}
	f.keys = append(f.keys, k)
}

type filterKey struct {
	key   string
	Type  ruleType
	rules []filterRule
}

type filterRule struct {
	Comparasion ComparasionType
	indices     []int
}

type ruleType uint8

const (
	ruleTypeInt ruleType = iota
	ruleTypeUint
	ruleTypeFloat
	ruleTypeString
	ruleTypeTime
)
