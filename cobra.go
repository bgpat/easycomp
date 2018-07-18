package easycomp

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type CobraCommand cobra.Command

func AddCobraCommand(rootCmd *cobra.Command, name string) *CommandLine {
	var c CommandLine
	rootCmd.AddCommand(&cobra.Command{
		Use:   name,
		Short: "Generate completion script",
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.ArgsLenAtDash() != 0 {
				s, err := dumpScript(os.Args[0], strings.Join(os.Args, " "))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Print(s)
				return
			}
			c.run(args)
		},
	})
	for _, cmd := range (*CobraCommand)(rootCmd).Children() {
		c.Append(cmd)
	}
	return &c
}

func (c *CobraCommand) Children() []Argument {
	args := []Argument{NewString(c.Name())}
	if c.IsRoot() {
		args = []Argument{}
	}
	for _, cmd := range c.Command().Commands() {
		args = append(
			args,
			NewCondition(
				Nest((*CobraCommand)(cmd), 1),
				func(words []string) bool {
					cur := current(words)
					if c.IsRoot() || c.Name() == cur {
						return true
					}
					for _, alias := range c.Command().Aliases {
						if alias == cur {
							return true
						}
					}
					return false
				},
			),
			Nest(NewPFlagSet(c.Command().Flags(), c), 1),
		)
	}
	for _, alias := range c.Command().Aliases {
		args = append(args, NewString(alias))
	}
	for _, arg := range c.Command().SuggestFor {
		args = append(args, NewString(arg))
	}
	for _, arg := range c.Command().ValidArgs {
		args = append(args, NewString(arg))
	}
	for p := c.Command(); p != nil; p = p.Parent() {
		args = append(
			args,
			Nest(NewPFlagSet(p.PersistentFlags(), c), 1),
		)
	}
	return args
}

func (c *CobraCommand) Complete(words []string) []string {
	args := Arguments(c.Children())
	return args.Complete(words)
}

func (c *CobraCommand) Match(words []string) bool {
	return true
}

func (c *CobraCommand) Command() *cobra.Command {
	return (*cobra.Command)(c)
}

func (c *CobraCommand) String() string {
	return c.Name()
}

func (c *CobraCommand) Name() string {
	return c.Command().Name()
}

func (c *CobraCommand) IsRoot() bool {
	return c.Command() == c.Command().Root()
}
