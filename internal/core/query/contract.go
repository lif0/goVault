package query

type Parser interface {
	Transition(in string) (out *Query, err error)
}

type Query struct {
	CommandID DBCommand
	Arguments []string
}
