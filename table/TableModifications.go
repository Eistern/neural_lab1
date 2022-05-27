package table

import "math"

func (table *DataTable) FillNaN(columnIndex int, value float64) {
	for _, rowValues := range table.Data {
		if math.IsNaN(rowValues[columnIndex]) {
			rowValues[columnIndex] = value
		}
	}
}

func (table *DataTable) RemoveLines(
	removePredicate func(lineValues []float64) bool) {
	for lineIndex := 0; lineIndex < len(table.Data); {
		if removePredicate(table.Data[lineIndex]) {
			table.removeLine(lineIndex)
		} else {
			lineIndex++
		}
	}
}

func (table *DataTable) removeLine(lineIndex int) {
	if lineIndex+1 != len(table.Data) {
		table.Data = append(table.Data[:lineIndex], table.Data[lineIndex+1:]...)
	} else {
		table.Data = table.Data[:lineIndex]
	}
}

func (table *DataTable) RemoveUnnamedColumns() {
	for columnNumber := 0; columnNumber < len(table.Labels); {
		if table.Labels[columnNumber] == "" {
			table.removeColumn(columnNumber)
		} else {
			columnNumber++
		}
	}
}

func (table *DataTable) removeColumn(columnNumber int) {
	if columnNumber+1 != len(table.Labels) {
		table.Labels = append(table.Labels[:columnNumber], table.Labels[columnNumber+1:]...)
		for rowIndex, rowValues := range table.Data {
			table.Data[rowIndex] = append(rowValues[:columnNumber], rowValues[columnNumber+1:]...)
		}
	} else {
		table.Labels = table.Labels[:columnNumber]
		for rowIndex, rowValues := range table.Data {
			table.Data[rowIndex] = rowValues[:columnNumber]
		}
	}
}

type MergeColumnIndexes struct {
	AcceptorColumn int
	DonorColumn    int
}

func (table *DataTable) MergeColumns(
	columnIndexes MergeColumnIndexes,
	mergeFunction func(acceptorColumnValue float64, donorColumnValue float64) float64) {
	for _, rowValues := range table.Data {
		rowValues[columnIndexes.AcceptorColumn] = mergeFunction(
			rowValues[columnIndexes.AcceptorColumn],
			rowValues[columnIndexes.DonorColumn])
	}
	table.removeColumn(columnIndexes.DonorColumn)
}
