package easycomp

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CommandLine struct {
	Arguments
	line  string
	point int
	cword int
}

func (c *CommandLine) Complete(words []string) []string {
	c.getenv()
	var partialWords []string
	if c.cword < 1 || len(words) < c.cword || len(c.line) < c.point {
		partialWords = words
	} else {
		partialWords = words[:c.cword]
		partialWord := words[c.cword]
		partialLine := c.line[:c.point]
		for i := len(partialWord); i >= 0; i-- {
			w := partialWord[:i]
			if w == partialLine[len(partialLine)-i:] {
				partialWords = append(partialWords, w)
			}
		}
	}
	return c.Arguments.Complete(partialWords[:c.cword+1])
}

func (c *CommandLine) getenv() {
	c.line = os.Getenv("COMP_LINE")
	c.point, _ = strconv.Atoi(os.Getenv("COMP_POINT"))
	c.cword, _ = strconv.Atoi(os.Getenv("COMP_CWORD"))
}

func (c *CommandLine) String() string {
	return ""
}

func (c *CommandLine) Set(_ string) error {
	for i, w := range os.Args {
		if w != "--" {
			continue
		}
		c.run(os.Args[i+1:])
		os.Exit(0)
	}
	s, err := dumpScript(os.Args[0], strings.Join(os.Args, " "))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(s)
	os.Exit(0)
	return nil
}

func (c *CommandLine) IsBoolFlag() bool {
	return true
}

func (c *CommandLine) Type() string {
	return ""
}

func (c *CommandLine) run(args []string) {
	debugCompletion(c, args)
	for _, s := range c.Complete(args) {
		fmt.Println(s)
	}
}
