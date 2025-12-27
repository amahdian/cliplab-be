package pg

import (
	"fmt"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type MetadataStg struct {
	db *gorm.DB
}

// Since the tables and columns do not change we can safely cache them
var cachedTableToColumnsMap map[string][]string

func NewMetadataStg(ses *ormSession) *MetadataStg {
	return &MetadataStg{db: ses.db}
}

const listTablesAndColumnsSql = `
SELECT t.table_name, c.column_name
FROM information_schema.tables t
JOIN information_schema.columns c ON t.table_name = c.table_name
WHERE t.table_schema = 'public'
  AND t.table_type = 'BASE TABLE'
ORDER BY t.table_name ASC, c.column_name ASC;
`

func (stg *MetadataStg) ListTablesAndColumns() (tableToColumnsMap map[string][]string, err error) {
	if cachedTableToColumnsMap == nil {
		type tableColumn struct {
			TableName  string
			ColumnName string
		}
		var tableColumns []tableColumn
		err = stg.db.Raw(listTablesAndColumnsSql).Scan(&tableColumns).Error

		if err == nil {
			cachedTableToColumnsMap = make(map[string][]string)
			for _, tableColumn := range tableColumns {
				if _, ok := cachedTableToColumnsMap[tableColumn.TableName]; !ok {
					cachedTableToColumnsMap[tableColumn.TableName] = make([]string, 0)
				}
				cachedTableToColumnsMap[tableColumn.TableName] = append(cachedTableToColumnsMap[tableColumn.TableName], tableColumn.ColumnName)
			}
		}
	}

	tableToColumnsMap = cachedTableToColumnsMap

	return
}

func (stg *MetadataStg) ListTables() (tableNames []string, err error) {
	tableToColumnsMap, err := stg.ListTablesAndColumns()

	if err != nil {
		tableNames = make([]string, 0)
		for tableName := range tableToColumnsMap {
			tableNames = append(tableNames, tableName)
		}
	}

	return
}

func (stg *MetadataStg) ListColumns(tableName string) (columnNames []string, err error) {
	tableToColumnsMap, err := stg.ListTablesAndColumns()

	if err != nil {
		if _, ok := tableToColumnsMap[tableName]; !ok {
			err = fmt.Errorf("Table '%s' does not exists.", tableName)
			return
		}

		columnNames = tableToColumnsMap[tableName]
	}

	return
}

func (stg *MetadataStg) RecordByValueExists(tableName, columnName string, value interface{}) (exists bool, err error) {
	exists, err = stg.anyRowsByValue(tableName, columnName, value)
	return
}

func (stg *MetadataStg) anyRowsByValue(tableName, columnName string, value interface{}) (exists bool, err error) {
	tableColumnsMap, err := stg.ListTablesAndColumns()
	if err != nil {
		return
	}
	if _, ok := tableColumnsMap[tableName]; !ok {
		err = fmt.Errorf("Table '%s' does not exists.", tableName)
		return
	}
	if !lo.Contains(tableColumnsMap[tableName], columnName) {
		err = fmt.Errorf("Column '%s' does not exists in table '%s'", columnName, tableName)
		return
	}

	res := stg.db.
		Table(tableName).
		Where(fmt.Sprintf("%s = ?", columnName), value).
		Select(columnName).
		Limit(1).
		Find(map[string]interface{}{})
	err = res.Error
	exists = res.RowsAffected > 0
	return
}
