package common

type SortOrder string

const (
	SortOrderAscending  SortOrder = "ASC"
	SortOrderDescending SortOrder = "DESC"

	DefaultSortOrder SortOrder = SortOrderAscending
)

type SearchCondition string

const (
	SearchConditionContains SearchCondition = "contains"
	SearchConditionEqual    SearchCondition = "eq"
	SearchConditionNotEqual SearchCondition = "neq"
)

const (
	DefaultPageSize = 100
)

type Pagination struct {
	// must be greater than 0. default: 100
	PageSize int `json:"pageSize" form:"pageSize" binding:"min=0"`

	// starts from 0
	Page int `json:"page" form:"page" binding:"min=0"`

	// field to sort the results by. if orderBy is empty, the order will be ignored.
	OrderBy string `json:"orderBy" form:"orderBy"`

	// sort order. default order is asc.
	// * asc - Ascending, from A to Z.
	// * desc - Descending, from Z to A.
	Order SortOrder `json:"order" form:"order" binding:"oneof=asc desc ASC DESC ''" enums:"asc,desc"`

	// for internal use only
	TotalCount int64 `swaggerignore:"true"`
}

func DefaultPagination() *Pagination {
	return &Pagination{
		Order:    DefaultSortOrder,
		PageSize: DefaultPageSize,
		Page:     0,
	}
}

var internalPagination = &Pagination{}

// InternalPagination returns a special pagination object that is only used in internal API calls.
// This special pagination signals the paginator to disable pagination for internal API calls and return all results.
func InternalPagination() *Pagination {
	return internalPagination
}

func NewPagination(orderBy string, order SortOrder) *Pagination {
	return &Pagination{
		OrderBy:  orderBy,
		Order:    order,
		PageSize: DefaultPageSize,
		Page:     0,
	}
}

type SearchParams struct {
	Filters []*FieldFilter `json:"filters" binding:"dive"`
	*Pagination
}

func DefaultSearchParams() *SearchParams {
	return &SearchParams{
		Filters:    make([]*FieldFilter, 0),
		Pagination: DefaultPagination(),
	}
}

type FieldFilter struct {
	// field name that should be used in search filter
	FieldName string `json:"fieldName" binding:"required"`

	// field condition.
	// * contains - Contains.
	// * eq - Equal.
	// * neq - Not equal.
	Condition SearchCondition `json:"condition" binding:"required,oneof=contains eq neq" enums:"contains eq neq"`

	// field value.
	Value string `json:"value"`
}
