package table

import "math"

type ColumnDistributionInfo struct {
	minValue     float64
	bucketSize   float64
	bucketCount  int
	bucketValues map[int]float64
}

func (table DataTable) ColumnDistribution(columnIndex int) ColumnDistributionInfo {
	var distributionInfo = ColumnDistributionInfo{}
	distributionInfo.minValue = table.ColumnMin(columnIndex)

	columnElementCount := table.Size() - table.ColumnNaNCount(columnIndex)
	distributionInfo.bucketCount = int(math.Floor(math.Sqrt(float64(columnElementCount))))

	bucketSplit := table.ColumnValueBucketSplit(columnIndex, distributionInfo.bucketCount)
	distributionInfo.bucketValues = map[int]float64{}
	for i := 0; i < distributionInfo.bucketCount; i++ {
		distributionInfo.bucketValues[i] = float64(bucketSplit[i]) / float64(columnElementCount)
	}

	return distributionInfo
}
