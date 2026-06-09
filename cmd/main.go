package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	dsox3024 "github.com/devicehub-go/keysight-dsox3024"
	dsoproto "github.com/devicehub-go/keysight-dsox3024/protocol"
	"github.com/devicehub-go/unicomm"
	"github.com/devicehub-go/unicomm/protocol/unicommtcp"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var channelColors = []color.RGBA{
	{R: 31, G: 119, B: 180, A: 255},
	{R: 255, G: 127, B: 14, A: 255},
	{R: 44, G: 160, B: 44, A: 255},
}

func plotWaveforms(filePath string, curves map[int]dsoproto.Waveform) error {
	p := plot.New()
	p.X.Label.Text = "Time (s)"
	p.Y.Label.Text = "Voltage (V)"
	p.Legend.Top = true
	p.Legend.Padding = vg.Points(-30)

	grid := plotter.NewGrid()
	lightGray := color.RGBA{R: 211, G: 211, B: 211, A: 255}
	grid.Vertical.Color = lightGray
	grid.Horizontal.Color = lightGray
	p.Add(grid)

	for i, ch := range []int{1, 2, 3} {
		waveform := curves[ch]
		if len(waveform.Timestamps) != len(waveform.Values) {
			return fmt.Errorf("ch%d: timestamps and values length mismatch", ch)
		}

		pts := make(plotter.XYs, len(waveform.Timestamps))
		for j := range pts {
			pts[j].X = waveform.Timestamps[j]
			pts[j].Y = waveform.Values[j]
		}

		line, err := plotter.NewLine(pts)
		if err != nil {
			return fmt.Errorf("ch%d: %w", ch, err)
		}
		line.LineStyle = draw.LineStyle{
			Color: channelColors[i],
			Width: vg.Points(1.2),
		}

		p.Add(line)
		p.Legend.Add(fmt.Sprintf("Ch%d", ch), line)
	}

	if err := p.Save(10*vg.Inch, 5*vg.Inch, filePath); err != nil {
		return fmt.Errorf("could not save plot: %w", err)
	}
	return nil
}

func saveWaveformTxt(filePath string, curves map[int]dsoproto.Waveform) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "%-12s %-12s %-12s %-12s\n", "Timestamp", "Ch1", "Ch2", "Ch3")

	n := len(curves[1].Values)
	for i := range n {
		_, err := fmt.Fprintf(file, "%-12f %-12f %-12f %-12f\n",
			curves[1].Timestamps[i],
			curves[1].Values[i],
			curves[2].Values[i],
			curves[3].Values[i],
		)
		if err != nil {
			return fmt.Errorf("failed to write row %d: %w", i, err)
		}
	}
	return nil
}

func main() {
	ip := flag.String("ip", "10.0.4.161", "Oscilloscope IP address")
	interval := flag.Int("interval", 10, "interval between acquisitions in seconds")
	flag.Parse()

	osc := dsox3024.New(unicomm.Options{
		Protocol: unicomm.TCP,
		TCP: unicommtcp.TCPOptions{
			Host:         *ip,
			Port:         5025,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
	})
	if err := osc.Connect(); err != nil {
		log.Fatal(err)
	}
	defer osc.Disconnect()

	if err := osc.SetWaveformPointsMode("NORMal"); err != nil {
		log.Fatal(err)
	}
	if err := osc.SetWaveformPoints(250); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll("./oscilloscope_data", 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		nowStr := time.Now().Format("2006_01_02_15_04_05")

		curves := make(map[int]dsoproto.Waveform)
		for _, ch := range []int{1, 2, 3} {
			waveform, err := osc.GetWaveform(ch)
			if err != nil {
				log.Fatal(err)
			}
			curves[ch] = *waveform
		}

		base := fmt.Sprintf("./oscilloscope_data/%s", nowStr)

		if err := plotWaveforms(base+".png", curves); err != nil {
			log.Fatal(err)
		}
		if err := saveWaveformTxt(base+".txt", curves); err != nil {
			log.Fatal(err)
		}
	}
}
