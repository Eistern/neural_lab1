package table

import (
	"fmt"
	"math"
)

func (table DataTable) PrintStatisticsElements() {
	for columnIndex, columnName := range table.Labels {
		columnEmptyCount := table.ColumnNaNCount(columnIndex)
		columnMin := table.ColumnMin(columnIndex)
		columnMax := table.ColumnMax(columnIndex)
		columnMean := table.ColumnMean(columnIndex)
		fmt.Printf("%s | EmptyCount - %d Min - %.3f Max - %.3f Mean - %.3f\n",
			columnName, columnEmptyCount, columnMin, columnMax, columnMean)
	}
}

func (table DataTable) ColumnDispersion(columnIndex int) float64 {
	valueMean := table.ColumnMean(columnIndex)
	valueDeviation := 0.0
	for _, rowValues := range table.Data {
		if !math.IsNaN(rowValues[columnIndex]) {
			valueDeviation += math.Pow(rowValues[columnIndex]-valueMean, 2)
		}
	}
	return math.Sqrt(valueDeviation)
}

func (table DataTable) ColumnMean(columnIndex int) float64 {
	valueCount := 0
	valueSum := 0.0
	for _, rowValues := range table.Data {
		if !math.IsNaN(rowValues[columnIndex]) {
			valueCount++
			valueSum += rowValues[columnIndex]
		}
	}
	return valueSum / float64(valueCount)
}

func (table DataTable) ColumnMin(columnIndex int) float64 {
	minimalValue := math.Inf(+1)
	for _, rowValues := range table.Data {
		if !math.IsNaN(rowValues[columnIndex]) {
			minimalValue = math.Min(minimalValue, rowValues[columnIndex])
		}
	}
	return minimalValue
}

func (table DataTable) ColumnMax(columnIndex int) float64 {
	maxValue := math.Inf(-1)
	for _, rowValues := range table.Data {
		if !math.IsNaN(rowValues[columnIndex]) {
			maxValue = math.Max(maxValue, rowValues[columnIndex])
		}
	}
	return maxValue
}

func (table DataTable) ColumnNaNCount(columnIndex int) int {
	nanCount := 0
	for _, rowValues := range table.Data {
		if math.IsNaN(rowValues[columnIndex]) {
			nanCount++
		}
	}
	return nanCount
}

func (table DataTable) ColumnsValueCount(columnIndices [2]int) map[[2]float64]int {
	valueCount := map[[2]float64]int{}
	for _, rowValues := range table.Data {
		elementValue := [2]float64{rowValues[columnIndices[0]], rowValues[columnIndices[1]]}
		valueCount[elementValue]++
	}
	return valueCount
}

func (table DataTable) ColumnValueCount(columnIndex int) map[float64]int {
	valueCount := map[float64]int{}
	for _, rowValues := range table.Data {
		elementValue := rowValues[columnIndex]
		valueCount[elementValue]++
	}
	return valueCount
}

func (table DataTable) ColumnValueBucketSplit(columnIndex, bucketCount int) map[int]int {
	columnMin := table.ColumnMin(columnIndex)
	bucketSize := int(table.ColumnMax(columnIndex)-columnMin) / bucketCount
	if bucketSize <= 0 {
		panic(fmt.Sprintf("Bucket split performed on a constant column: %d", columnIndex))
	}

	var valueBucketCount map[int]int
	for _, rowValues := range table.Data {
		if !math.IsNaN(rowValues[columnIndex]) {
			elementBucketIndex := int(rowValues[columnIndex]-columnMin) / bucketSize
			valueBucketCount[elementBucketIndex]++
		}
	}
	return valueBucketCount
}
