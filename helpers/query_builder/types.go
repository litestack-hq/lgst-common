package query_builder

type TableSortOrder string

func (value *TableSortOrder) IsValid() bool {
	if value == nil {
		return false
	}
	return *value == "ASC" || *value == "DESC"
}

type PaginationQueryInput struct {
	InitialQuery string
	Limit        int
	NextCursor   string
	Sort         struct {
		Field       string
		Order       TableSortOrder
		CursorValue any
	}
	Search struct {
		Query  string
		Fields []string
	}
}
