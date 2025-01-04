//go:generate mockgen -destination ./../../../mocks/core/query/contract.go -package ${GOPACKAGE}_mock . Parser
package query

type Parser interface {
	Transition(in string) (out *Query, err error)
}

type Query struct {
	CommandID DBCommand
	Arguments []string
}
