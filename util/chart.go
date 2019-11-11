package util

import (
	"fmt"
	"github.com/chenjiandongx/go-echarts/charts"
	"log"
	"os"
)

/**
 * 折线图
 */
func ChartLine(name string, xAxis []interface{}, yAxis []interface{}) {
	showName := "line-" + name
	bar := charts.NewLine()
	bar.SetGlobalOptions(charts.TitleOpts{Title: showName})
	bar.AddXAxis(xAxis).
		AddYAxis(showName, yAxis)

	f, err := os.Create(fmt.Sprintf("./charts/%s.html", showName))
	if err != nil {
		log.Println(err)
	}

	err = bar.Render(f)
	bar.Overlap()
	if err != nil {
		log.Print(err)
	}
}

/**
 * kline + 点位
 */
func ChartOver(name string, xAxis []interface{}, yAxis [][]interface{}, lines ...[]interface{}) {
	bar:= charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: "Kline_Target"},
		charts.XAxisOpts{SplitNumber: 20},
		charts.YAxisOpts{Scale: true},
		charts.DataZoomOpts{XAxisIndex: []int{0}, Start: 50, End: 100},
	)
	bar.AddXAxis(xAxis)

	kline := charts.NewKLine()
	kline.AddXAxis(xAxis).AddYAxis("kline", yAxis)
	bar.Overlap(kline)

	for _,value := range lines {
		line := charts.NewLine()
		line.SetGlobalOptions(charts.TitleOpts{Title: "Target"})
		line.AddXAxis(xAxis).AddYAxis("", value)
		bar.Overlap(line)
	}

	f, err := os.Create(fmt.Sprintf("./charts/%s.html", name))
	if err != nil {
		log.Println(err)
	}
	err = bar.Render(f)
	if err != nil {
		log.Print(err)
	}
}

/**
 * 3D 图形
 */
func Chart3D(name string, axis [][3]interface{}) {
	showName := "3D-" + name
	scatter3d := charts.NewScatter3D()
	scatter3d.SetGlobalOptions(
		charts.TitleOpts{Title: showName},
		charts.VisualMapOpts{
			Calculable: true,
		},
		charts.Grid3DOpts{BoxDepth: 80, BoxWidth: 200},
	)
	scatter3d.AddZAxis(showName, axis)
	f, err := os.Create(fmt.Sprintf("./charts/%s.html", showName))
	if err != nil {
		log.Println(err)
	}

	err = scatter3d.Render(f)
	if err != nil {
		log.Print(err)
	}
}
