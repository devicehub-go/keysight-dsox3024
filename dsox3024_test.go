package dsox3024_test

import (
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
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	})
	if err := osc.Connect(); err != nil {
		t.Fatalf("error on connect: %v", err)
	}
	defer osc.Disconnect()

	ini := time.Now()
	waveform, err := osc.GeTriggeredtWaveform(1, 5*time.Second)
	if err != nil {
		t.Fatalf("error on get waveform: %v", err)
	}
	t.Log(time.Since(ini))
	t.Log(len(waveform.Timestamps))
	t.Log(len(waveform.Values))
}
