package qfl

import (
	"time"
)

type RuleType uint8

const (
	RuleTypeInt RuleType = iota
	RuleTypeUint
	RuleTypeFloat
	RuleTypeString
	RuleTypeTime
)

type ComparasionType uint8

const (
	ComparasionInvalid ComparasionType = iota
	ComparasionEquals
	ComparasionLessThan
	ComparasionMoreThan
	ComparasionLessOrEqual
	ComparasionMoreOrEqual
	ComparasionLike
	ComparasionIs
)

type FilterRule[T comparable] struct {
	Comparasion ComparasionType
	Values      []T
}

// Filter is a specialized data structure that stores rules for a given key. It
// only supports some primitive data types. All get and set functions should be
// the exactly same except for the type it's manipulating, this is on purporse
// to avoid casting costs and type safety.
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
		if f.keys[i].key == key && f.keys[i].Type == RuleTypeInt {
			return getGeneric(f.keys[i], f.intVals)
		}
	}

	return nil
}

func (f *Filter) GetUint(key string) []FilterRule[uint] {
	for i := range f.keys {
		if f.keys[i].key == key && f.keys[i].Type == RuleTypeUint {
			return getGeneric(f.keys[i], f.uintVals)
		}
	}

	return nil
}

func (f *Filter) GetFloat(key string) []FilterRule[float64] {
	for i := range f.keys {
		if f.keys[i].key == key && f.keys[i].Type == RuleTypeFloat {
			return getGeneric(f.keys[i], f.floatVals)
		}
	}

	return nil
}

func (f *Filter) GetString(key string) []FilterRule[string] {
	for i := range f.keys {
		if f.keys[i].key == key && f.keys[i].Type == RuleTypeString {
			return getGeneric(f.keys[i], f.stringVals)
		}
	}

	return nil
}

func (f *Filter) GetTime(key string) []FilterRule[time.Time] {
	for i := range f.keys {
		if f.keys[i].key == key && f.keys[i].Type == RuleTypeTime {
			return getGeneric(f.keys[i], f.timeVals)
		}
	}

	return nil
}

func getGeneric[T comparable](key filterKey, vals []T) []FilterRule[T] {
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
	f.appendRule(key, indices, comparasion, RuleTypeInt)
}

func (f *Filter) AddUint(key string, values []uint, comparasion ComparasionType) {
	start := len(f.uintVals)
	f.uintVals = append(f.uintVals, values...)
	end := len(f.uintVals)

	indices := generateSequence(start, end)
	f.appendRule(key, indices, comparasion, RuleTypeUint)
}

func (f *Filter) AddFloat(key string, values []float64, comparasion ComparasionType) {
	start := len(f.floatVals)
	f.floatVals = append(f.floatVals, values...)
	end := len(f.floatVals)

	indices := generateSequence(start, end)
	f.appendRule(key, indices, comparasion, RuleTypeFloat)
}

func (f *Filter) AddString(key string, values []string, comparasion ComparasionType) {
	start := len(f.stringVals)
	f.stringVals = append(f.stringVals, values...)
	end := len(f.stringVals)

	indices := generateSequence(start, end)
	f.appendRule(key, indices, comparasion, RuleTypeString)
}

func (f *Filter) AddTime(key string, values []time.Time, comparasion ComparasionType) {
	start := len(f.timeVals)
	f.timeVals = append(f.timeVals, values...)
	end := len(f.timeVals)

	indices := generateSequence(start, end)
	f.appendRule(key, indices, comparasion, RuleTypeTime)
}

func (f *Filter) appendRule(key string, indices []int, comparasion ComparasionType, ruleType RuleType) {
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
	Type  RuleType
	rules []filterRule
}

type filterRule struct {
	Comparasion ComparasionType
	indices     []int
}
