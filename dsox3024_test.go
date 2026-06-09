package dsox3024_test

import (
	"fmt"
	"testing"
	"time"

	dsox3024 "github.com/devicehub-go/keysight-dsox3024"
	"github.com/devicehub-go/unicomm"
	"github.com/devicehub-go/unicomm/protocol/unicommtcp"
)

func TestOscilloscope(t *testing.T) {
	osc := dsox3024.New(unicomm.Options{
		Protocol: unicomm.TCP,
		TCP: unicommtcp.TCPOptions{
			Host:         "10.0.4.161",
			Port:         5025,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
	})
	if err := osc.Connect(); err != nil {
		t.Fatal(err)
	}
	defer osc.Disconnect()

	osc.SetWaveformPointsMode("NORMal")
	osc.SetWaveformPoints(250)
	points, _ := osc.Query(":WAVEform:POINts?")
	fmt.Println(string(points))
	waveform, err := osc.GetWaveform(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(waveform.Timestamps))
	t.Log(len(waveform.Values))
}
