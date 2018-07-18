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
	return true
}

func (s *FlagSet) String() string {
	return strings.Join(s.Complete([]string{""}), ", ")
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
		NewString(f.String()),
		NewCondition(
			Nest((*FlagSet)(f.flagSet), 1),
			func(words []string) bool {
				return strings.HasPrefix(current(words), f.String()+"=")
			},
		),
		NewCondition(
			Nest((*FlagSet)(f.flagSet), 2),
			func(words []string) bool {
				return current(words) == f.String()
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

func (f *Flag) String() string {
	return "-" + f.Name
}

func (f *Flag) IsBoolFlag() bool {
	bf, ok := f.Flag.Value.(boolFlag)
	return ok && bf.IsBoolFlag()
}

func AddFlag(flagSet *flag.FlagSet, name, desc string) *CommandLine {
	var c CommandLine
	flagSet.Var(&c, name, desc)
	c.Append(Nest((*FlagSet)(flagSet), 1))
	return &c
}
