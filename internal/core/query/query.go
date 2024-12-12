package query

func NewQuery(commandID DBCommand) Query {
	return Query{
		CommandID: commandID,
	}
}
