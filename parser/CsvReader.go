package parser

import (
	"dataStatsProject/table"
	"encoding/csv"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

func Parse(filename string) table.DataTable {
	fileDescriptor, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func(fileDescriptor *os.File) {
		err := fileDescriptor.Close()
		if err != nil {
			panic(err)
		}
	}(fileDescriptor)

	csvReader := csv.NewReader(fileDescriptor)
	csvReader.Comma = ';'
	dataLines, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	return toDataTable(dataLines)
}

func toDataTable(stringData [][]string) table.DataTable {
	result := table.DataTable{Labels: stringData[0]}

	dataRowCount := len(stringData) - 2
	dataColumnCount := len(stringData[0])

	result.Data = make([][]float64, dataRowCount)

	nilValue := math.NaN()
	lock := sync.WaitGroup{}
	for rowIndex, rowValues := range stringData[2:] {
		result.Data[rowIndex] = make([]float64, dataColumnCount)

		for columnIndex, columnValue := range rowValues {

			lock.Add(1)
			go func(value string, rowIndex int, columnNumber int) {
				defer lock.Done()
				floatValue, err := strconv.ParseFloat(toAmericanFloatRepresentation(value), 32)
				if err != nil {
					floatValue = nilValue
				}
				result.Data[rowIndex][columnNumber] = floatValue
			}(columnValue, rowIndex, columnIndex)

		}
	}
	lock.Wait()
	return result
}

func toAmericanFloatRepresentation(europeanRepresentation string) string {
	return strings.ReplaceAll(europeanRepresentation, ",", ".")
}
