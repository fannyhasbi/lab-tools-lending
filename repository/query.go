package repository

type QueryResult struct {
	Result interface{}
	Error  error
}

type ToolQuery interface {
	GetTool() QueryResult
}
