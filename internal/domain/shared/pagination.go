package shared

// Page holds pagination parameters for repository queries.
type Page struct {
	Number int // 1-based page number
	Size   int // items per page (max 500)
}

const (
	DefaultPageSize = 20
	MaxPageSize     = 500
)

// Cursor-based alternative for high-volume time-series queries.
type Cursor struct {
	After  string
	Before string
	Limit  int
}

// PageResult wraps a list result with pagination metadata.
type PageResult[T any] struct {
	Items      []T
	TotalItems int64
	Page       int
	PageSize   int
	TotalPages int
}

func NewPage(page, size int) Page {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = DefaultPageSize
	}
	if size > MaxPageSize {
		size = MaxPageSize
	}
	return Page{Number: page, Size: size}
}

func (p Page) Offset() int {
	p = p.Normalized()
	return (p.Number - 1) * p.Size
}

// Normalized enforces repository-safe pagination bounds even when callers
// construct Page values directly instead of using NewPage.
func (p Page) Normalized() Page {
	return NewPage(p.Number, p.Size)
}

func NewPageResult[T any](items []T, total int64, p Page) PageResult[T] {
	p = p.Normalized()
	totalPages := int(total) / p.Size
	if int(total)%p.Size != 0 {
		totalPages++
	}
	return PageResult[T]{
		Items:      items,
		TotalItems: total,
		Page:       p.Number,
		PageSize:   p.Size,
		TotalPages: totalPages,
	}
}

// SortOrder defines ascending or descending ordering.
type SortOrder string

const (
	SortAsc  SortOrder = "ASC"
	SortDesc SortOrder = "DESC"
)

// Sort specifies a sort field and direction.
type Sort struct {
	Field string
	Order SortOrder
}
