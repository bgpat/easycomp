package easycomp

import (
	"errors"
)

var (
	notImplementedError = errors.New("not implemented")
)

func current(words []string) string {
	if len(words) > 0 {
		return words[0]
	}
	return ""
}

func setSingle(dst Argument, src []Argument) error {
	if len(src) == 0 {
		return errors.New("not enough arguments")
	}
	if len(src) > 1 {
		return errors.New("too many arguments")
	}
	dst = src[0]
	return nil
}
