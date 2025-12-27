package pg

const (
	//maxMssqlParams  = 2100
	maxPgParams = 65535
	//maxOracleParams = 1000

	maxAllowedParams = maxPgParams
)

func CalcBestBatchSize[T any](_ []T) int {
	batchSize := maxAllowedParams / ColumnCount[T]()
	return batchSize
}

func ColumnCount[T any]() int {
	var model T
	schema := getGormSchema(&model)
	return len(schema.DBNames)
}
