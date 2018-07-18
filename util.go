package easycomp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	debugOutput string
	debugDepth  = 5

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

type Tree struct {
	Type     string
	Value    string
	Match    bool
	Children []*Tree
	Words    string
}

func debugCompletion(a Argument, words []string) {
	if debugOutput == "" || debugDepth <= 0 {
		return
	}
	b, err := json.Marshal(toTree(a, words, debugDepth))
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("out.json", b, 0644)
}

func toTree(arg Argument, words []string, depth int) *Tree {
	if arg == nil || len(words) < 1 {
		return nil
	}
	if depth < 0 {
		return &Tree{Type: "too deep"}
	}
	c := make([]*Tree, 0)
	if _, ok := arg.(*Nester); ok {
		words = words[1:]
	}
	for _, a := range arg.Children() {
		tree := toTree(a, words, depth-1)
		if tree == nil {
			continue
		}
		c = append(c, tree)
	}
	return &Tree{
		Type:     fmt.Sprintf("%T", arg),
		Value:    fmt.Sprintf("%v", arg),
		Match:    arg.Match(words),
		Children: c,
		Words:    strings.Join(words, ", "),
	}
	return nil
}
