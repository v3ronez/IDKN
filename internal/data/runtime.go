package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsValue := fmt.Sprintf("%d mins", r)

	return []byte(strconv.Quote(jsValue)), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	v, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	parts := strings.Split(v, " ")
	if len(parts) > 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}
	*r = Runtime(i)
	return nil
}
