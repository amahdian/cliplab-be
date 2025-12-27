package paged

import (
	"fmt"
	"strings"
)

type PaginatedRequest struct {
	Page int    `form:"page,default=0"`
	Size int    `form:"size,default=20"`
	Sort string `form:"sort"`
}

func GetSortDetails(sort string) (string, int, error) {
	var sortOrder int
	sortDetails := strings.Split(sort, ",")
	if len(sortDetails) != 2 {
		return "", 0, fmt.Errorf("invalid Sort Options")
	}
	field := sortDetails[0]
	order := sortDetails[1]
	if order == "asc" {
		sortOrder = 1
	} else if order == "desc" {
		sortOrder = -1
	} else {
		return "", 0, fmt.Errorf("invalid sortOrder value. Use 'ASC' or 'DESC'")
	}
	return field, sortOrder, nil
}
