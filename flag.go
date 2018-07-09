package easycomp

import (
	"flag"
	"strings"
)

type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

type FlagSet struct {
	Arguments
}

func NewFlagSet(flagSet *flag.FlagSet) *FlagSet {
	if flagSet == nil {
		return nil
	}
	var arg FlagSet
	flagSet.VisitAll(func(f *flag.Flag) {
		fa := NewFlag(f)
		arg.Append(fa)
		arg.Append(NewCondition(Nest(&arg, 1), func(words []string) bool {
			return strings.HasPrefix(current(words), fa.String()+"=")
		}))
		arg.Append(NewCondition(Nest(&arg, 2), func(words []string) bool {
			return current(words) == fa.String()
		}))
	})
	return &arg
}

type Flag struct {
	*flag.Flag
}

func NewFlag(f *flag.Flag) *Flag {
	if f == nil {
		return nil
	}
	return &Flag{Flag: f}
}

func (f *Flag) Get() []Argument {
	return []Argument{}
}

func (f *Flag) Set(_ []Argument) error {
	return notImplementedError
}

func (f *Flag) Name() string {
	return f.String()
}

func (f *Flag) Complete(words []string) []string {
	if len(words) == 1 {
		return []string{f.String()}
	}
	return nil
}

func (f *Flag) Match(words []string) bool {
	return strings.HasPrefix(f.String(), current(words))
}

func (f *Flag) String() string {
	return "-" + f.Flag.Name
}

func (f *Flag) IsBoolFlag() bool {
	bf, ok := f.Flag.Value.(boolFlag)
	return ok && bf.IsBoolFlag()
}
