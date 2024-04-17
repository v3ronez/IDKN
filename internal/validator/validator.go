package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("/^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/")
)

type Validator struct {
	Errors map[string]string
}

// helper to create a new validator
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key string, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key string, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func PermittdValue[T comparable](value T, permitted ...T) bool {
	for i := range permitted {
		if value == permitted[i] {
			return true
		}
	}
	return false
}

func Matches(s string, rgx *regexp.Regexp) bool {
	return rgx.MatchString(s)
}

func Unique[T comparable](v []T) bool {
	uniqueValue := make(map[T]bool)
	for _, i := range v {
		uniqueValue[i] = true
	}
	return len(uniqueValue) == len(v)
}
