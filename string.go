package easycomp

import (
	"strings"
)

type String string

func NewString(s string) *String {
	d := String(s)
	return &d
}

func (s *String) Complete(words []string) []string {
	return []string{s.String()}
}

func (s *String) Match(words []string) bool {
	return strings.HasPrefix(s.String(), current(words))
}

func (s *String) String() string {
	return string(*s)
}
