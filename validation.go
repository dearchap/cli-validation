package validation

import (
	"fmt"
)

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

type Float interface {
	~float32 | ~float64
}

type Ordered interface {
	Integer | Float
}

// ConditionOrError ir a helper function to make writing
// validation functions much easier
func ConditionOrError(cond bool, err error) error {
	if cond {
		return nil
	}
	return err
}

// ValidationChainAll allows one to chain a sequence of validation
// functions to construct a single validation function. All the
// individual validations must pass for the validation to succeed
func ValidationChainAll[T any](fns ...func(T) error) func(T) error {
	return func(v T) error {
		for _, fn := range fns {
			if err := fn(v); err != nil {
				return err
			}
		}
		return nil
	}
}

// ValidationChainAny allows one to chain a sequence of validation
// functions to construct a single validation function. Atleast one
// of the individual validations must pass for the validation to succeed
func ValidationChainAny[T any](fns ...func(T) error) func(T) error {
	return func(v T) error {
		var errs []error
		for _, fn := range fns {
			if err := fn(v); err == nil {
				return nil
			} else {
				errs = append(errs, err)
			}
		}
		return fmt.Errorf("%+v", errs)
	}
}

// Min means that the value to be checked needs to be atleast(and including)
// the checked value
func Min[T Ordered](c T) func(T) error {
	return func(v T) error {
		return ConditionOrError(v >= c, fmt.Errorf("%v is not less than %v", v, c))
	}
}

// Max means that the value to be checked needs to be atmost(and including)
// the checked value
func Max[T Ordered](c T) func(T) error {
	return func(v T) error {
		return ConditionOrError(v <= c, fmt.Errorf("%v is not greater than %v", v, c))
	}
}

// Max means that the value to be checked needs to be atmost(and including)
// the checked value
func RangeInclusive[T Ordered](a, b T) func(T) error {
	return ValidationChainAll[T](Min[T](a), Max[T](b))
}
