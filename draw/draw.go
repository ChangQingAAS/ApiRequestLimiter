package main

import (
	"github.com/ChangQingAAS/ApiRequestLimiter/csvUtils"
	"log"
	"math/rand"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	rand.Seed(int64(0))

	p := plot.New()

	p.Title.Text = "请求数和token展示"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Tokens"

	err := plotutil.AddLinePoints(p,
		"numRequests", handleLog(1),
		"reserveTokens", handleLog(2),
		"currTokens", handleLog(3))
	if err != nil {
		log.Fatal(err)
	}

	if err = p.Save(10*vg.Inch, 10*vg.Inch, "./draw/points.png"); err != nil {
		log.Fatal(err)
	}
}

func handleLog(temp int) plotter.XYs {
	allData := csvUtils.ReadCsv("./csvUtils/log.csv")
	points := make(plotter.XYs, len(allData))
	i := 0
	for _, item := range allData {
		y, _ := strconv.ParseFloat(item[0], 64)
		x, _ := strconv.ParseFloat(item[temp], 64)
		points[i].X = y
		points[i].Y = x
		i++
	}
	return points
}
