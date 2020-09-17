package flagtype

import (
	"fmt"
)

// FlagArr custom type to parse flag arr
type FlagArr []string

// String impl flag.Value::String
func (a FlagArr) String() string {
	return fmt.Sprint([]string(a))
}

// Set impl flag.Value::Set
func (a *FlagArr) Set(val string) error {
	*a = append(*a, val)
	return nil
}
