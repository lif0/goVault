package query

type DBCommand uint8

const (
	SET = DBCommand(1)
	GET = DBCommand(2)
	DEL = DBCommand(3)
)
