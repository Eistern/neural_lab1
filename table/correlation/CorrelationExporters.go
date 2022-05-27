package correlation

import (
	"dataStatsProject/table"
	"encoding/csv"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"log"
	"math"
	"os"
	"strconv"

	"gonum.org/v1/plot/palette"
)

func ExportCorrelationTableCsv(outputFile string, data table.DataTable) {
	outputCsv, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer func(outputCsv *os.File) {
		err := outputCsv.Close()
		if err != nil {
			panic(err)
		}
	}(outputCsv)

	csvWriter := csv.NewWriter(outputCsv)
	csvWriter.UseCRLF = false
	correlationTable := table.CorrelationTable(data,
		func(correlationValue float64) string {
			return strconv.FormatFloat(math.Abs(correlationValue), 'f', 3, 64)
		})
	err = csvWriter.WriteAll(correlationTable)
	if err != nil {
		panic(err)
	}
}

func ExportCorrelationTableHeatmap(outputFile string, data table.DataTable) {
	correlationValues := table.CorrelationTable(data, func(correlationValue float64) float64 {
		return correlationValue
	})

	pal := palette.Heat(12, 1)
	heatMap := plotter.NewHeatMap(toXYZGrid(correlationValues), pal)

	p := plot.New()
	p.Title.Text = "Heat map"

	p.X.Tick.Marker = labelTicks{data.Labels, 10}
	p.Y.Tick.Marker = labelTicks{data.Labels, 10}

	p.Add(heatMap)

	// Create a legend.
	l := plot.NewLegend()

	img := vgimg.New(1024, 1024)
	dc := draw.New(img)

	l.Top = true
	// Calculate the width of the legend.
	r := l.Rectangle(dc)
	legendWidth := r.Max.X - r.Min.X
	l.YOffs = -p.Title.TextStyle.FontExtents().Height // Adjust the legend down a little.

	l.Draw(dc)
	dc = draw.Crop(dc, 0, -legendWidth-vg.Millimeter, 0, 0) // Make space for the legend.
	p.Draw(dc)
	w, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Panic(err)
	}
	png := vgimg.PngCanvas{Canvas: img}
	if _, err = png.WriteTo(w); err != nil {
		log.Panic(err)
	}
}

type labelTicks struct {
	labelValues []string
	shift       int
}

func (labelTicker labelTicks) Ticks(min, max float64) []plot.Tick {
	var t []plot.Tick
	for i := math.Trunc(min); i <= max; i++ {
		indexValue := int(i) / labelTicker.shift
		if int(i)%labelTicker.shift == 0 {
			if indexValue >= 0 && indexValue < len(labelTicker.labelValues) {
				t = append(t, plot.Tick{Value: i, Label: labelTicker.labelValues[indexValue]})
			}
		}
	}
	return t
}

type deciGrid struct {
	mat.Matrix
}

func (g deciGrid) Dims() (c, r int) {
	r, c = g.Matrix.Dims()
	return c, r
}
func (g deciGrid) Z(c, r int) float64 { return g.Matrix.At(r, c) }
func (g deciGrid) X(c int) float64 {
	_, n := g.Matrix.Dims()
	if c < 0 || c >= n {
		panic("index out of range")
	}
	return 10 * float64(c)
}
func (g deciGrid) Y(r int) float64 {
	m, _ := g.Matrix.Dims()
	if r < 0 || r >= m {
		panic("index out of range")
	}
	return 10 * float64(r)
}

func toXYZGrid(values [][]float64) plotter.GridXYZ {
	flatten := make([]float64, len(values)*len(values))
	for i := 0; i < len(flatten); i++ {
		flatten[i] = values[i/len(values)][i%len(values)]
	}

	matrix := mat.NewDense(len(values), len(values), flatten)

	return deciGrid{matrix}
}
