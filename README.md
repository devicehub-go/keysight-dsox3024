# keysight-dsox3024

A Go library for communicating with and controlling the **Keysight DSOX3024** oscilloscope over TCP (SCPI/LAN interface).

## Installation

```bash
go get github.com/devicehub-go/keysight-dsox3024
```

## Quick Start

```go
import (
    dsox3024 "github.com/devicehub-go/keysight-dsox3024"
    "github.com/devicehub-go/unicomm"
    "github.com/devicehub-go/unicomm/protocol/unicommtcp"
    "time"
)

osc := dsox3024.New(unicomm.Options{
    Protocol: unicomm.TCP,
    TCP: unicommtcp.TCPOptions{
        Host:         "192.168.1.100",
        Port:         5025,
        ReadTimeout:  1 * time.Second,
        WriteTimeout: 1 * time.Second,
    },
})

if err := osc.Connect(); err != nil {
    log.Fatal(err)
}
defer osc.Disconnect()
```

## API Reference

### Connection

| Method | Description |
|---|---|
| `Connect() error` | Connects to the device and verifies identity via `*IDN?` |
| `Disconnect() error` | Closes the connection |
| `IsConnected() bool` | Returns true if the device is currently connected |

### Communication

| Method | Description |
|---|---|
| `Write(message string) error` | Sends a SCPI command to the device |
| `WriteSequence(messages []string) error` | Sends multiple SCPI commands in order |
| `Query(message string) ([]byte, error)` | Sends a command and reads a newline-terminated response |
| `QueryByteSequence(message string) ([]byte, error)` | Sends a command and reads an IEEE 488.2 binary block response |

### Waveform Acquisition

| Method | Description |
|---|---|
| `GetWaveform(channel int) (*Waveform, error)` | Acquires and converts waveform data from the given channel (1–4) |
| `GetPreamble(channel int) (Preamble, error)` | Retrieves waveform scaling metadata (time/voltage increments, origins, references) |
| `SetWaveformPoints(points int) error` | Sets the number of waveform points to transfer |
| `SetWaveformPointsMode(mode string) error` | Sets the points mode: `NORMal`, `RAW`, or `MAXimum` |
| `ByteToFloatArray(payload []byte) ([]float64, error)` | Converts a raw SWAP 64-bit little-endian byte payload to float64 values |

The `Waveform` struct returned by `GetWaveform` contains:

```go
type Waveform struct {
    Timestamps []float64 // seconds
    Values     []float64 // volts
}
```

### Vertical (Voltage) Scale

| Method | Description |
|---|---|
| `SetVerticalScale(channel int, units float64) error` | Sets volts-per-division for the given channel |
| `GetVerticalScale(channel int) (float64, error)` | Gets current volts-per-division for the given channel |

### Timebase

| Method | Description |
|---|---|
| `SetTimebaseScale(units float64) error` | Sets seconds-per-division |
| `SetTimebaseRange(fullScale float64) error` | Sets the full-scale horizontal time window in seconds |
| `GetTimebaseRange() (float64, error)` | Gets the current full-scale horizontal time window |
| `SetTimebasePosition(position float64) error` | Sets the delay between the trigger event and the reference point |
| `GetTimebasePosition() (float64, error)` | Gets the current timebase delay in seconds |
| `SetTimebaseReference(location string) error` | Sets the horizontal reference position: `LEFT`, `RIGHt`, or `CENTer` |
| `GetTimebaseReference() (string, error)` | Gets the current horizontal reference position |

### Measurements & Statistics

| Method | Description |
|---|---|
| `GetVoltageAmplitude(channel int) (float64, error)` | Returns the difference between the top and base voltage of the waveform |
| `GetVoltageAverage(channel int) (float64, error)` | Returns the average voltage over the displayed waveform |
| `GetStatistics() ([]Statistics, error)` | Enables statistics display and returns all current measurement results |

The `Statistics` struct contains `Label`, `Current`, `Minimum`, `Maximum`, `Mean`, `Std`, and `Count` fields.

## Example: Capturing Three Channels

```go
osc.SetWaveformPointsMode("NORMal")
osc.SetWaveformPoints(250)

for _, ch := range []int{1, 2, 3} {
    waveform, err := osc.GetWaveform(ch)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Ch%d: %d samples\n", ch, len(waveform.Values))
}
```

## Dependencies

- [`github.com/devicehub-go/unicomm`](https://github.com/devicehub-go/unicomm) — transport abstraction (TCP, Serial)
- [`gonum.org/v1/plot`](https://pkg.go.dev/gonum.org/v1/plot) — waveform plotting (used by the example application)