package table

import "math"

func (table DataTable) ColumnGainRatio(columnIndex int, targetColumnIndices [2]int) float64 {
	return table.columnInformationGain(columnIndex, targetColumnIndices) / table.columnIntrinsicInformation(columnIndex)
}

func (table DataTable) columnIntrinsicInformation(columnIndex int) float64 {
	return table.columnEntropy(columnIndex)
}

func (table DataTable) columnInformationGain(columnIndex int, targetColumnIndices [2]int) float64 {
	baseEntropy := table.columnsEntropy(targetColumnIndices)

	columnValueCount := table.ColumnValueCount(columnIndex)
	for columnValue := range columnValueCount {
		baseEntropy -= table.valueInformationGain(columnIndex, columnValue, targetColumnIndices)
	}
	return baseEntropy
}

func (table DataTable) valueInformationGain(valueColumnIndex int, value float64, targetColumnIndices [2]int) float64 {
	modifiedTable := table.Clone()
	modifiedTable.RemoveLines(func(lineValues []float64) bool {
		return lineValues[valueColumnIndex] != value
	})

	modifiedEntropy := modifiedTable.columnsEntropy(targetColumnIndices)
	return modifiedEntropy
}

func (table DataTable) columnsEntropy(targetColumnIndices [2]int) float64 {
	values := table.ColumnsValueCount(targetColumnIndices)
	entropy := 0.0
	for _, count := range values {
		valueProbability := float64(count) / float64(table.Size())
		entropy -= valueProbability * math.Log2(valueProbability)
	}
	return entropy
}

func (table DataTable) columnEntropy(targetColumnIndex int) float64 {
	values := table.ColumnValueCount(targetColumnIndex)
	entropy := 0.0
	for _, count := range values {
		valueProbability := float64(count) / float64(table.Size())
		entropy -= valueProbability * math.Log2(valueProbability)
	}
	return entropy
}
