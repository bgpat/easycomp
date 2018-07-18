package easycomp

type Argument interface {
	Children() []Argument
	Complete([]string) []string
	Match([]string) bool
	String() string
}

type Arguments []Argument

func (a *Arguments) Children() []Argument {
	return *a
}

func (a *Arguments) Append(arg Argument) {
	*a = append(*a, arg)
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
	return len(*a) > 0
}

func (a *Arguments) String() string {
	return ""
}

type Nester struct {
	Argument
}

func Nest(arg Argument, depth int) Argument {
	if depth == 0 {
		return arg
	}
	return &Nester{Argument: Nest(arg, depth-1)}
}

func (n *Nester) Complete(words []string) []string {
	if len(words) < 1 {
		return nil
	}
	return n.Argument.Complete(words[1:])
}

func (n *Nester) Match(words []string) bool {
	return len(words) >= 1
}

type Condition struct {
	Argument
	match func([]string) bool
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
