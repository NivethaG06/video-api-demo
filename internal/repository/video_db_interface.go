package repository

type DB interface {
	InsertData(name string, age int) error
	FetchPaginatedData(page, pageSize int) ([]string, error)
	CreateIndexPost() error
	CreateIndexTag() error
}
