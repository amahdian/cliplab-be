package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Metadata map[string]string

func (m *Metadata) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil // Store nil if the map is nil
	}
	j, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}
	return string(j), nil // GORM will store this as TEXT/JSONB
}

func (m *Metadata) Scan(src interface{}) error {
	if src == nil {
		*m = nil // Set to nil if the database value is NULL
		return nil
	}

	var sourceBytes []byte
	switch s := src.(type) {
	case []byte:
		sourceBytes = s
	case string:
		sourceBytes = []byte(s)
	default:
		return errors.New("incompatible type for Metadata: expected []byte or string")
	}

	if len(sourceBytes) == 0 {
		*m = make(Metadata) // Initialize as an empty map if database value is empty string/bytes
		return nil
	}

	// Make sure the target map is initialized before unmarshaling
	// If *m is nil (from previous scan), json.Unmarshal needs a non-nil target.
	if *m == nil {
		*m = make(Metadata)
	}

	if err := json.Unmarshal(sourceBytes, m); err != nil {
		return fmt.Errorf("failed to unmarshal metadata from JSON: %w", err)
	}
	return nil
}
