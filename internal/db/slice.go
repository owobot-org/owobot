package db

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	return strings.Join(s, "\x1F"), nil
}

func (s *StringSlice) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("incompatible type for StringSlice")
	}
	*s = splitOptions(str)
	return nil
}

func splitOptions(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\x1F")
}
