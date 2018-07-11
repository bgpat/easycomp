package easycomp

import (
	"strings"
)

type String string

func NewString(s string) *String {
	a := String(s)
	return &a
}

func (s *String) Children() []Argument {
	return nil
}

func (s *String) Complete(words []string) []string {
	if len(words) == 1 {
		return []string{s.String()}
	}
	return nil
}

func (s *String) Match(words []string) bool {
	return strings.HasPrefix(s.String(), current(words))
}

func (s *String) String() string {
	return string(*s)
}
