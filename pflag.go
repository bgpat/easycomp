package easycomp

import (
	"strings"

	"github.com/spf13/pflag"
)

type PFlagSet struct {
	*pflag.FlagSet
	cmd *CobraCommand
}

func NewPFlagSet(s *pflag.FlagSet, cmd *CobraCommand) *PFlagSet {
	return &PFlagSet{FlagSet: s, cmd: cmd}
}

func (s *PFlagSet) Children() []Argument {
	flags := make([]Argument, 0, s.NFlag())
	s.VisitAll(func(f *pflag.Flag) {
		flags = append(flags, NewPFlag(f, s.cmd))
	})
	return flags
}

func (s *PFlagSet) Complete(words []string) []string {
	args := Arguments(s.Children())
	return args.Complete(words)
}

func (s *PFlagSet) Match(words []string) bool {
	return s.NFlag() > 0
}

func (s *PFlagSet) String() string {
	return "[flag] " + s.cmd.String()
}

type PFlag struct {
	*pflag.Flag
	cmd *CobraCommand
}

func NewPFlag(f *pflag.Flag, c *CobraCommand) *PFlag {
	if c.Command().DisableFlagParsing {
		return nil
	}
	return &PFlag{f, c}
}

func (f *PFlag) Children() []Argument {
	c := Arguments(f.cmd.Children())
	args := []Argument{
		NewString(f.String()),
		NewCondition(
			&c,
			func(words []string) bool {
				return strings.HasPrefix(current(words), f.String()+"=")
			},
		),
		NewCondition(
			Nest(&c, 1),
			func(words []string) bool {
				return current(words) == f.String()
			},
		),
	}
	if f.Shorthand != "" {
		args = append(
			args,
			NewString("-"+f.Shorthand),
			NewCondition(
				&c,
				func(words []string) bool {
					return strings.HasPrefix(current(words), "-"+f.Shorthand)
				},
			),
			NewCondition(
				Nest(&c, 1),
				func(words []string) bool {
					return current(words) == "-"+f.Shorthand
				},
			),
		)
	}
	return args
}

func (f *PFlag) Complete(words []string) []string {
	args := Arguments(f.Children())
	return args.Complete(words)
}

func (f *PFlag) Match(words []string) bool {
	return strings.HasPrefix(current(words), "-")
}

func (f *PFlag) String() string {
	return "--" + f.Name
}
