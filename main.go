package main

import (
	csv "dataStatsProject/parser"
	"dataStatsProject/table"
	"dataStatsProject/table/correlation"
	"fmt"
	"math"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage: <path_to_csv>")
		return
	}

	dataTable := csv.Parse(os.Args[1])
	dataTable.RemoveUnnamedColumns()

	dataTable.MergeColumns(table.MergeColumnIndexes{
		AcceptorColumn: len(dataTable.Labels) - 2,
		DonorColumn:    len(dataTable.Labels) - 1},
		func(acceptorColumnValue float64, donorColumnValue float64) float64 {
			if math.IsNaN(acceptorColumnValue) {
				return donorColumnValue * 1000
			}
			return acceptorColumnValue
		})

	dataTable.RemoveLines(
		func(lineValues []float64) bool {
			return math.IsNaN(lineValues[len(dataTable.Labels)-1]) && math.IsNaN(lineValues[len(dataTable.Labels)-2])
		})

	for columnIndex, columnName := range dataTable.Labels {
		fmt.Printf("%s | %.2f\n", columnName, dataTable.ColumnGainRatio(columnIndex, [2]int{len(dataTable.Labels) - 2, len(dataTable.Labels) - 1}))
	}

	for labelIndex := range dataTable.Labels {
		columnMean := dataTable.ColumnMean(labelIndex)
		dataTable.FillNaN(labelIndex, columnMean)
	}

	correlation.ExportCorrelationTableCsv("correlation.csv", dataTable)
	correlation.ExportCorrelationTableHeatmap("correlation_heatmap.png", dataTable)
}
