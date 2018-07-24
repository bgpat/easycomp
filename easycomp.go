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

type Transform struct {
	Argument
	fn func([]string) []string
}

func NewTransform(arg Argument, fn func([]string) []string) *Transform {
	return &Transform{Argument: arg, fn: fn}
}

func Nest(arg Argument, depth int) Argument {
	if depth == 0 {
		return arg
	}
	return NewTransform(
		Nest(arg, depth-1),
		func(words []string) []string {
			if len(words) < 1 {
				return []string{}
			}
			return words[1:]
		},
	)
}

func (t *Transform) Transform(words []string) []string {
	return t.fn(words)
}

func (t *Transform) Complete(words []string) []string {
	return t.Argument.Complete(t.Transform(words))
}

func (t *Transform) Match(words []string) bool {
	return len(t.Transform(words)) >= 1
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
