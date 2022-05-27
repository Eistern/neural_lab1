package table

import (
	"fmt"
)

type DataTable struct {
	Labels []string
	Data   [][]float64
}

func (table DataTable) Clone() DataTable {
	clonedLabels := make([]string, len(table.Labels))
	copy(clonedLabels, table.Labels)

	clonedData := make([][]float64, len(table.Data))
	for clonedRowIndex := range clonedData {
		clonedData[clonedRowIndex] = make([]float64, len(table.Data[clonedRowIndex]))
		copy(clonedData[clonedRowIndex], table.Data[clonedRowIndex])
	}

	return DataTable{
		Labels: clonedLabels,
		Data:   nil,
	}
}

func (table DataTable) Size() int {
	return len(table.Data)
}

func (table *DataTable) Print() {
	for _, label := range table.Labels {
		print(label, "\t")
	}
	print("\n")

	for _, row := range table.Data {
		for _, element := range row {
			fmt.Printf("%.2f\t", element)
		}
		print("\n")
	}
}
