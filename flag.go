package easycomp

import (
	"flag"
	"strings"
)

type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

type FlagSet flag.FlagSet

func (s *FlagSet) Children() []Argument {
	flags := make([]Argument, 0)
	(*flag.FlagSet)(s).VisitAll(func(f *flag.Flag) {
		flags = append(flags, NewFlag(f, (*flag.FlagSet)(s)))
	})
	return flags
}

func (s *FlagSet) Complete(words []string) []string {
	args := Arguments(s.Children())
	return args.Complete(words)
}

func (s *FlagSet) Match(words []string) bool {
	return len(words) > 0
}

type Flag struct {
	*flag.Flag
	flagSet *flag.FlagSet
}

func NewFlag(f *flag.Flag, s *flag.FlagSet) *Flag {
	return &Flag{f, s}
}

func (f *Flag) Children() []Argument {
	return []Argument{
		NewString("-" + f.Flag.Name),
		NewCondition(
			Nest((*FlagSet)(f.flagSet), 1),
			func(words []string) bool {
				return strings.HasPrefix(current(words), "-"+f.Name+"=")
			},
		),
		NewCondition(
			Nest((*FlagSet)(f.flagSet), 2),
			func(words []string) bool {
				return current(words) == "-"+f.Name
			},
		),
	}
}

func (f *Flag) Complete(words []string) []string {
	args := Arguments(f.Children())
	return args.Complete(words)
}

func (f *Flag) Match(words []string) bool {
	return len(words) > 0
}

func (f *Flag) IsBoolFlag() bool {
	bf, ok := f.Flag.Value.(boolFlag)
	return ok && bf.IsBoolFlag()
}
