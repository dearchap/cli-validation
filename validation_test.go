package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testcase[T any] struct {
	name        string
	f           func(T) error
	input       T
	errExpected bool
}

func testWithInput[T any](t *testing.T, tcases []testcase[T]) {
	r := require.New(t)
	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.f(tc.input)
			if tc.errExpected {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}
func TestMin(t *testing.T) {

	inputs := []testcase[int]{
		{
			name:  "min lower limit",
			f:     Min[int](10),
			input: 10,
		},
		{
			name:  "min normal",
			f:     Min[int](10),
			input: 11,
		},
		{
			name:        "min failure",
			f:           Min[int](10),
			input:       8,
			errExpected: true,
		},
	}

	testWithInput[int](t, inputs)
}
func TestMax(t *testing.T) {

	inputs := []testcase[float64]{
		{
			name:  "max upper limit",
			f:     Max[float64](10.0),
			input: 10.0,
		},
		{
			name:  "max normal",
			f:     Max[float64](10.0),
			input: 9.0,
		},
		{
			name:        "max failure",
			f:           Max[float64](10.0),
			input:       10.00001,
			errExpected: true,
		},
	}

	testWithInput[float64](t, inputs)
}
func TestRange(t *testing.T) {

	inputs := []testcase[uint64]{
		{
			name:  "lower limit",
			f:     RangeInclusive[uint64](10, 16),
			input: 10,
		},
		{
			name:  "upper limit",
			f:     RangeInclusive[uint64](10, 16),
			input: 16,
		},
		{
			name:        "lower limit failure",
			f:           RangeInclusive[uint64](10, 16),
			input:       9,
			errExpected: true,
		},
		{
			name:        "upper limit failure",
			f:           RangeInclusive[uint64](10, 16),
			input:       17,
			errExpected: true,
		},
		{
			name:  "normal",
			f:     RangeInclusive[uint64](10, 16),
			input: 12,
		},
	}

	testWithInput[uint64](t, inputs)
}
func TestDisjointRange(t *testing.T) {

	inputs := []testcase[int16]{
		{
			name:  "lower limit lower range",
			f:     ValidationChainAny[int16](RangeInclusive[int16](10, 16), RangeInclusive[int16](56, 67)),
			input: 10,
		},
		{
			name:  "upper limit lower range",
			f:     ValidationChainAny[int16](RangeInclusive[int16](10, 16), RangeInclusive[int16](56, 67)),
			input: 16,
		},
		{
			name:  "lower limit upper range",
			f:     ValidationChainAny[int16](RangeInclusive[int16](10, 16), RangeInclusive[int16](56, 67)),
			input: 56,
		},
		{
			name:  "upper limit upper range",
			f:     ValidationChainAny[int16](RangeInclusive[int16](10, 16), RangeInclusive[int16](56, 67)),
			input: 67,
		},
		{
			name:  "normal lower range",
			f:     ValidationChainAny[int16](RangeInclusive[int16](10, 16), RangeInclusive[int16](56, 67)),
			input: 13,
		},
		{
			name:  "normal upper range",
			f:     ValidationChainAny[int16](RangeInclusive[int16](10, 16), RangeInclusive[int16](56, 67)),
			input: 60,
		},
		{
			name:        "failure below lower range",
			f:           ValidationChainAny[int16](RangeInclusive[int16](10, 16), RangeInclusive[int16](56, 67)),
			input:       1,
			errExpected: true,
		},
		{
			name:        "failure in between lower and upper range",
			f:           ValidationChainAny[int16](RangeInclusive[int16](10, 16), RangeInclusive[int16](56, 67)),
			input:       20,
			errExpected: true,
		},
		{
			name:        "failure above upper range",
			f:           ValidationChainAny[int16](RangeInclusive[int16](10, 16), RangeInclusive[int16](56, 67)),
			input:       70,
			errExpected: true,
		},
	}

	testWithInput[int16](t, inputs)
}
func TestEnum(t *testing.T) {

	inputs := []testcase[string]{
		{
			name:  "valid value for enum - 1",
			f:     Enum[string]("hello", "bar", "foo"),
			input: "hello",
		},
		{
			name:  "valid value for enum - 2",
			f:     Enum[string]("hello", "bar", "foo"),
			input: "bar",
		},
		{
			name:        "invalid value for enum",
			f:           Enum[string]("hello", "bar", "foo"),
			input:       "helloee",
			errExpected: true,
		},
	}

	testWithInput[string](t, inputs)
}

func TestRegex(t *testing.T) {
	type myString string

	inputs := []testcase[myString]{
		{
			name:  "valid value regex match - 1",
			f:     Regex[myString]("foo[1-7].*y"),
			input: "foo1y",
		},
		{
			name:  "valid value regex match - 2",
			f:     Regex[myString]("foo[1-7].*y"),
			input: "foo1ttthsyy",
		},
		{
			name:        "invalid value regex match - 1",
			f:           Regex[myString]("foo[1-7].*y"),
			input:       "fooy",
			errExpected: true,
		},
		{
			name:        "invalid value regex match - 2",
			f:           Regex[myString]("foo[1-7].*y"),
			input:       "fooOy",
			errExpected: true,
		},
	}

	testWithInput[myString](t, inputs)
}

func TestSliceValidator(t *testing.T) {

	inputs := []testcase[[]uint32]{
		{
			name:  "valid values slices",
			f:     SliceValidator[uint32](Min[uint32](7)),
			input: []uint32{9, 10, 18, 14},
		},
		{
			name:        "invalid value slice - 1",
			f:           SliceValidator[uint32](Min[uint32](7)),
			input:       []uint32{9, 10, 6, 14},
			errExpected: true,
		},
	}

	testWithInput[[]uint32](t, inputs)
}
