package jsonutil

import (
	"encoding/json"
	"fmt"
)

type StringOrNumber string

func (s *StringOrNumber) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*s = StringOrNumber(str)
		return nil
	}

	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		*s = StringOrNumber(num.String())
		return nil
	}

	return fmt.Errorf("data must be string or number")
}
