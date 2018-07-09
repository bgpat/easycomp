package easycomp

type Argument interface {
	Complete([]string) []string
	Match([]string) bool
}

type Arguments []Argument

type Nester struct {
	arg Argument
}

type Condition struct {
	Argument
	match func([]string) bool
}

func (a *Arguments) Complete(words []string) []string {
	reply := make([]string, 0)
	for _, arg := range *a {
		if arg.Match(words) {
			reply = append(reply, arg.Complete(words)...)
		}
	}
	return reply
}

func (a *Arguments) Match(_ []string) bool {
	return true
}

func (a *Arguments) Append(arg Argument) {
	*a = append(*a, arg)
}

func NewNester(arg Argument) *Nester {
	return &Nester{arg: arg}
}

func Nest(arg Argument, depth int) Argument {
	if depth == 0 {
		return arg
	}
	return NewNester(Nest(arg, depth-1))
}

func (n *Nester) Complete(words []string) []string {
	if len(words) < 1 {
		return nil
	}
	return n.arg.Complete(words[1:])
}

func (n *Nester) Match(words []string) bool {
	return len(words) >= 1
}

func NewCondition(arg Argument, match func([]string) bool) *Condition {
	return &Condition{
		Argument: arg,
		match:    match,
	}
}

func (c *Condition) Match(words []string) bool {
	return c.match(words)
}
