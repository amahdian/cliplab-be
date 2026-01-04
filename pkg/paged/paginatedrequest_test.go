package paged

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
)

func TestPaginatedRequestBinding(t *testing.T) {
	tests := []struct {
		name          string
		queryString   string
		expectedPage  int
		expectedSize  int
		expectedSort  string
		expectedError bool
	}{
		{
			name:          "Default values",
			queryString:   "",
			expectedPage:  0,
			expectedSize:  20,
			expectedSort:  "",
			expectedError: false,
		},
		{
			name:          "Custom values",
			queryString:   "page=2&size=10&sort=name,asc",
			expectedPage:  2,
			expectedSize:  10,
			expectedSort:  "name,asc",
			expectedError: false,
		},
		{
			name:          "Invalid page",
			queryString:   "page=invalid&size=10",
			expectedError: true,
		},
		{
			name:          "Invalid size",
			queryString:   "page=0&size=invalid",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin router
			gin.SetMode(gin.TestMode)
			router := gin.New()

			router.GET("/test", func(c *gin.Context) {
				var req PaginatedRequest
				if err := c.ShouldBindQuery(&req); (err != nil) != tt.expectedError {
					t.Errorf("ShouldBindQuery() error = %v, expectedError %v", err, tt.expectedError)
					return
				}

				if !tt.expectedError {
					assert.Equal(t, tt.expectedPage, req.Page)
					assert.Equal(t, tt.expectedSize, req.Size)
					assert.Equal(t, tt.expectedSort, req.Sort)
				}

				c.Status(http.StatusOK)
			})

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test?"+tt.queryString, nil)
			router.ServeHTTP(w, req)
		})
	}
}

func TestGetSortDetails(t *testing.T) {
	tests := []struct {
		name          string
		sort          string
		expectedField string
		expectedOrder int
		expectedError bool
	}{
		{
			name:          "Valid ascending sort",
			sort:          "name,asc",
			expectedField: "name",
			expectedOrder: 1,
			expectedError: false,
		},
		{
			name:          "Valid descending sort",
			sort:          "age,desc",
			expectedField: "age",
			expectedOrder: -1,
			expectedError: false,
		},
		{
			name:          "Invalid sort format - missing comma",
			sort:          "name_asc",
			expectedField: "",
			expectedOrder: 0,
			expectedError: true,
		},
		{
			name:          "Invalid sort format - empty string",
			sort:          "",
			expectedField: "",
			expectedOrder: 0,
			expectedError: true,
		},
		{
			name:          "Invalid sort order",
			sort:          "name,invalid",
			expectedField: "",
			expectedOrder: 0,
			expectedError: true,
		},
		{
			name:          "Extra comma in sort string",
			sort:          "name,asc,extra",
			expectedField: "",
			expectedOrder: 0,
			expectedError: true,
		},
		{
			name:          "Case sensitivity test for asc",
			sort:          "name,ASC",
			expectedField: "",
			expectedOrder: 0,
			expectedError: true,
		},
		{
			name:          "Case sensitivity test for desc",
			sort:          "name,DESC",
			expectedField: "",
			expectedOrder: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, order, err := GetSortDetails(tt.sort)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedField, field)
				assert.Equal(t, tt.expectedOrder, order)
			}
		})
	}
}

func TestPaginatedRequestValidation(t *testing.T) {
	tests := []struct {
		name          string
		request       PaginatedRequest
		expectedValid bool
	}{
		{
			name: "Valid request with defaults",
			request: PaginatedRequest{
				Page: 0,
				Size: 20,
			},
			expectedValid: true,
		},
		{
			name: "Valid request with custom values",
			request: PaginatedRequest{
				Page: 1,
				Size: 10,
				Sort: "name,asc",
			},
			expectedValid: true,
		},
		{
			name: "Invalid negative page",
			request: PaginatedRequest{
				Page: -1,
				Size: 20,
			},
			expectedValid: true,
		},
		{
			name: "Invalid zero size",
			request: PaginatedRequest{
				Page: 0,
				Size: 0,
			},
			expectedValid: true,
		},
		{
			name: "Invalid negative size",
			request: PaginatedRequest{
				Page: 0,
				Size: -1,
			},
			expectedValid: true,
		},
	}

	validate := validator.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.expectedValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
