package repository

type QueryResult struct {
	Result interface{}
	Error  error
}

type ToolsQuery interface {
	GetTools() QueryResult
}
