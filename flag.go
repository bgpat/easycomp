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
		flags = append(flags, NewFlag(f, s))
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
	flagSet *FlagSet
}

func NewFlag(f *flag.Flag, s *FlagSet) *Flag {
	return &Flag{f, s}
}

func (f *Flag) Children() []Argument {
	return []Argument{
		NewString(f.String()),
		NewCondition(
			Nest(NewFlagValue(f), 1),
			func(words []string) bool {
				return current(words) == f.String() && !f.IsBoolFlag()
			},
		),
		NewCondition(
			NewTransform(
				NewFlagValue(f),
				func(words []string) []string {
					cur := strings.TrimPrefix(current(words), f.String()+"=")
					return append([]string{cur}, words[1:]...)
				},
			),
			func(words []string) bool {
				return strings.HasPrefix(current(words), f.String()+"=")
			},
		),
	}
}

func (f *Flag) Complete(words []string) []string {
	args := Arguments(f.Children())
	return args.Complete(words)
}

func (f *Flag) Match(words []string) bool {
	return true
}

func (f *Flag) String() string {
	return "-" + f.Name
}

func (f *Flag) IsBoolFlag() bool {
	bf, ok := f.Flag.Value.(boolFlag)
	return ok && bf.IsBoolFlag()
}

func NewFlagValue(f *Flag) *FlagValue {
	return &FlagValue{
		Value: f.Flag.Value,
		flag:  f,
	}
}

func AddFlag(flagSet *flag.FlagSet, name, desc string) *CommandLine {
	var c CommandLine
	flagSet.Var(&c, name, desc)
	c.Append(Nest((*FlagSet)(flagSet), 1))
	return &c
}

type FlagValue struct {
	flag.Value
	flag *Flag
}

func (v *FlagValue) Children() []Argument {
	children := []Argument{Nest(v.flag.flagSet, 1)}
	if v.flag.IsBoolFlag() {
		children = append(children, NewString("true"), NewString("false"))
	}
	return children
}

func (v *FlagValue) Complete(words []string) []string {
	args := Arguments(v.Children())
	return args.Complete(words)
}

func (v *FlagValue) Match(words []string) bool {
	return true
}
