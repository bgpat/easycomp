package easycomp

import (
	"flag"
	"fmt"
	"log"
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
	if err := c.getenv(); err != nil {
		log.Fatal(err)
	}
	var partialWords []string
	if len(words) < c.cword || len(c.line) < c.point {
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
	return Nest(&c.Arguments, 1).Complete(partialWords[:c.cword+1])
}

func (c *CommandLine) getenv() (err error) {
	c.line = os.Getenv("COMP_LINE")
	c.point, err = strconv.Atoi(os.Getenv("COMP_POINT"))
	if err != nil {
		return
	}
	c.cword, err = strconv.Atoi(os.Getenv("COMP_CWORD"))
	return
}

func AddFlag(flagSet *flag.FlagSet, name, desc string) *CommandLine {
	var c CommandLine
	flagSet.Var(&c, name, desc)
	c.Append((*FlagSet)(flagSet))
	return &c
}

func (c *CommandLine) String() string {
	return ""
}

func (c *CommandLine) Set(_ string) error {
	for i, w := range os.Args {
		if w != "--" {
			continue
		}
		for _, s := range (*CommandLine)(c).Complete(os.Args[i+1:]) {
			fmt.Println(s)
		}
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
