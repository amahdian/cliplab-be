package paged

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSortDetail(t *testing.T) {
	tests := []struct {
		name     string
		sort     string
		expected SortDetail
	}{
		{
			name: "Empty sort string",
			sort: "",
			expected: SortDetail{
				Sorted:   false,
				Unsorted: true,
				Empty:    false,
			},
		},
		{
			name: "Non-empty sort string",
			sort: "name:asc",
			expected: SortDetail{
				Sorted:   true,
				Unsorted: false,
				Empty:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSortDetail(tt.sort)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewPaginatedResponse(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		dataCount  int
		totalCount int64
		page       int
		limit      int
		sort       string
		expected   *PaginatedResponse
	}{
		{
			name:       "Empty result set",
			data:       []string{},
			dataCount:  0,
			totalCount: 0,
			page:       0,
			limit:      10,
			sort:       "",
			expected: &PaginatedResponse{
				Content:          []string{},
				Empty:            true,
				Number:           0,
				NumberOfElements: 0,
				Size:             10,
				TotalElements:    0,
				TotalPages:       0,
				First:            true,
				Last:             true,
				Sort: SortDetail{
					Sorted:   false,
					Unsorted: true,
					Empty:    false,
				},
				Pageable: PaginationDetail{
					Offset:     0,
					PageNumber: 0,
					PageSize:   10,
					Paged:      true,
					Unpaged:    false,
					Sort: SortDetail{
						Sorted:   false,
						Unsorted: true,
						Empty:    false,
					},
				},
			},
		},
		{
			name:       "First page with data",
			data:       []string{"item1", "item2"},
			dataCount:  2,
			totalCount: 5,
			page:       0,
			limit:      2,
			sort:       "name:asc",
			expected: &PaginatedResponse{
				Content:          []string{"item1", "item2"},
				Empty:            false,
				Number:           0,
				NumberOfElements: 2,
				Size:             2,
				TotalElements:    5,
				TotalPages:       3,
				First:            true,
				Last:             false,
				Sort: SortDetail{
					Sorted:   true,
					Unsorted: false,
					Empty:    false,
				},
				Pageable: PaginationDetail{
					Offset:     0,
					PageNumber: 0,
					PageSize:   2,
					Paged:      true,
					Unpaged:    false,
					Sort: SortDetail{
						Sorted:   true,
						Unsorted: false,
						Empty:    false,
					},
				},
			},
		},
		{
			name:       "Last page with data",
			data:       []string{"item5"},
			dataCount:  1,
			totalCount: 5,
			page:       2,
			limit:      2,
			sort:       "name:asc",
			expected: &PaginatedResponse{
				Content:          []string{"item5"},
				Empty:            false,
				Number:           2,
				NumberOfElements: 1,
				Size:             2,
				TotalElements:    5,
				TotalPages:       3,
				First:            false,
				Last:             true,
				Sort: SortDetail{
					Sorted:   true,
					Unsorted: false,
					Empty:    false,
				},
				Pageable: PaginationDetail{
					Offset:     4,
					PageNumber: 2,
					PageSize:   2,
					Paged:      true,
					Unpaged:    false,
					Sort: SortDetail{
						Sorted:   true,
						Unsorted: false,
						Empty:    false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPaginatedResponse(tt.data, tt.dataCount, int(tt.totalCount), tt.page, tt.limit, tt.sort)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPaginatedResponseJSON(t *testing.T) {
	// Test JSON marshaling/unmarshaling
	response := NewPaginatedResponse(
		[]string{"item1", "item2"},
		2,
		5,
		0,
		2,
		"name:asc",
	)

	// Marshal to JSON
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)

	// Unmarshal back
	var unmarshaled PaginatedResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	// Check JSON tags work correctly
	assert.Equal(t, response.TotalElements, unmarshaled.TotalElements)
	assert.Equal(t, response.Size, unmarshaled.Size)
	assert.Equal(t, response.Number, unmarshaled.Number)
	assert.Equal(t, response.Sort.Sorted, unmarshaled.Sort.Sorted)
	assert.Equal(t, response.Pageable.PageSize, unmarshaled.Pageable.PageSize)
}
