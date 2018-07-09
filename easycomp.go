package easycomp

type Argument interface {
	Get() []Argument
	Set([]Argument) error
	Name() string
	Complete([]string) []string
	Match([]string) bool
}

type Arguments []Argument

func (a *Arguments) Get() []Argument {
	return *a
}

func (a *Arguments) Set(args []Argument) error {
	*a = args
	return nil
}

func (a *Arguments) Name() string {
	return ""
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

type Nester struct {
	arg Argument
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

func (n *Nester) Get() []Argument {
	return []Argument{n.arg}
}

func (n *Nester) Set(args []Argument) error {
	return setSingle(n.arg, args)
}

func (n *Nester) Name() string {
	return ""
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

func (c *Condition) Get() []Argument {
	return []Argument{c.Argument}
}

func (c *Condition) Set(args []Argument) error {
	return setSingle(c.Argument, args)
}

func (c *Condition) Name() string {
	return ""
}

func (c *Condition) Match(words []string) bool {
	return c.match(words)
}
