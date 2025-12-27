package pg

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var schemaCache = &sync.Map{}

type gormScope func(*gorm.DB) *gorm.DB

func getTablePrefix(tableAlias ...string) string {
	prefix := ""
	if len(tableAlias) > 0 {
		prefix = tableAlias[0] + "."
	}
	return prefix
}

func getGormSchema(model interface{}) *schema.Schema {
	s, err := schema.Parse(&model, schemaCache, schema.NamingStrategy{})
	if err != nil {
		panic("failed to parse table schema")
	}
	return s
}

func withAlias(table schema.Tabler, alias string) string {
	return fmt.Sprintf("%s AS %s", table.TableName(), alias)
}
