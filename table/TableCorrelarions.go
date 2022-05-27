package table

import (
	"math"
)

func CorrelationTable[T float64 | string](
	table DataTable,
	converter func(correlationValue float64) T) [][]T {
	tableSize := len(table.Labels)
	correlationTable := make([][]T, tableSize)
	for i := 0; i < tableSize; i++ {
		correlationTable[i] = make([]T, tableSize)
	}

	for firstColumnIndex := 0; firstColumnIndex < tableSize; firstColumnIndex++ {
		for secondColumnIndex := 0; secondColumnIndex <= firstColumnIndex; secondColumnIndex++ {
			columnsCorrelation := table.columnsCorrelation(firstColumnIndex, secondColumnIndex)
			correlationValue := converter(columnsCorrelation)
			correlationTable[firstColumnIndex][secondColumnIndex] = correlationValue
			correlationTable[secondColumnIndex][firstColumnIndex] = correlationValue
		}
	}
	return correlationTable
}

func (table DataTable) columnsCorrelation(firstColumnIndex, secondColumnIndex int) float64 {
	firstColumnMean := table.ColumnMean(firstColumnIndex)
	secondColumnMean := table.ColumnMean(secondColumnIndex)
	firstColumnDispersion := table.ColumnDispersion(firstColumnIndex)
	secondColumnDispersion := table.ColumnDispersion(secondColumnIndex)

	combinedColumnDeviation := 0.0
	for _, rowValues := range table.Data {
		if !math.IsNaN(rowValues[firstColumnIndex]) && !math.IsNaN(rowValues[secondColumnIndex]) {
			combinedColumnDeviation += (rowValues[firstColumnIndex] - firstColumnMean) * (rowValues[secondColumnIndex] - secondColumnMean)
		}
	}
	return combinedColumnDeviation / (firstColumnDispersion * secondColumnDispersion)
}
