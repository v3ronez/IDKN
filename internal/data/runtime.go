package data

import (
	"fmt"
	"strconv"
)

type Runtinme int32

func (r Runtinme) MarshalJSON() ([]byte, error) {
	jsValue := fmt.Sprintf("%d mins", r)

	return []byte(strconv.Quote(jsValue)), nil
}
