package pagination

type Pagination struct {
	Last       bool  `json:"last"`
	Cursor     int64 `json:"cursor"`
	NextCursor int64 `json:"next_cursor"`
	Size       int   `json:"size"`
	TotalCount int64 `json:"total_count"`
}

func Paginate(totalCount int64, cursor int64, size, defaultSize int, nextCursor int64) Pagination {
	if defaultSize == 0 {
		defaultSize = 10
	}

	if size == 0 {
		size = defaultSize
	}

	startIndex := cursor

	nextIndex := startIndex + int64(size)
	last := nextIndex >= totalCount

	if last {
		nextCursor = cursor
	}

	return Pagination{
		Cursor:     cursor,
		NextCursor: nextCursor,
		Size:       size,
		TotalCount: totalCount,
		Last:       last,
	}
}
