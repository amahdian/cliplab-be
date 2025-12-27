package storage

type MetadataStorage interface {
	ListTablesAndColumns() (tableToColumnsMap map[string][]string, err error)
	ListTables() (tableNames []string, err error)
	ListColumns(tableName string) (columnNames []string, err error)
	RecordByValueExists(tableName, columnName string, value interface{}) (exists bool, err error)
}
