package protocol

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Preamble struct {
	Format     int
	Type       string
	Points     int
	Count      int
	xIncrement float64
	xOrigin    float64
	xReference float64
	yIncrement float64
	yOrigin    float64
	yReference float64
}

type Statistics struct {
	Label   string
	Current float64
	Minimum float64
	Maximum float64
	Mean    float64
	Std     float64
	Count   float64
}

type Waveform struct {
	Timestamps []float64
	Values     []float64
}

// Gets the preamble for the desired channel
func (d *DSOX3024) GetPreamble(channel int) (Preamble, error) {
	var preamble Preamble

	response, err := d.Query(":WAVeform:PREamble?")
	if err != nil {
		return preamble, err
	}

	preambleData := strings.Split(string(response), ",")
	if len(preambleData) < 10 {
		return preamble, fmt.Errorf("invalid preamble, got %s", response)
	}

	preamble.Format, _ = strconv.Atoi(preambleData[0])
	preamble.Type = preambleData[1]
	preamble.Points, _ = strconv.Atoi(preambleData[2])
	preamble.Count, _ = strconv.Atoi(preambleData[3])
	preamble.xIncrement, _ = strconv.ParseFloat(preambleData[4], 64)
	preamble.xOrigin, _ = strconv.ParseFloat(preambleData[5], 64)
	preamble.xReference, _ = strconv.ParseFloat(preambleData[6], 64)
	preamble.yIncrement, _ = strconv.ParseFloat(preambleData[7], 64)
	preamble.yOrigin, _ = strconv.ParseFloat(preambleData[8], 64)
	preamble.yReference, _ = strconv.ParseFloat(preambleData[9], 64)

	return preamble, nil
}

// Get statistics from the measured signals
func (d *DSOX3024) GetStatistics() ([]Statistics, error) {
	statistics := make([]Statistics, 0)

	if err := d.WriteSequence([]string{
		"SYSTem:MENU MEASure",
		"MEASure:STATistics:DISPlay ON",
		"MEASure:STATistics ON",
	}); err != nil {
		return nil, err
	}
	response, err := d.Query("MEASure:RESults?")
	if err != nil {
		return nil, err
	}

	cols := 7
	values := strings.Split(string(response), ",")
	if len(values)%cols != 0 {
		return nil, fmt.Errorf("unexpected statistics response, got %s", response)
	}

	for row := 0; row < len(values)/cols; row++ {
		current, _ := strconv.ParseFloat(values[row*cols+1], 64)
		minimum, _ := strconv.ParseFloat(values[row*cols+2], 64)
		maximum, _ := strconv.ParseFloat(values[row*cols+3], 64)
		mean, _ := strconv.ParseFloat(values[row*cols+4], 64)
		std, _ := strconv.ParseFloat(values[row*cols+5], 64)
		count, _ := strconv.ParseFloat(values[row*cols+6], 64)
		statistics = append(statistics, Statistics{
			Label:   values[row*cols],
			Current: current,
			Minimum: minimum,
			Maximum: maximum,
			Mean:    mean,
			Std:     std,
			Count:   count,
		})
	}

	return statistics, nil
}

// Calculates the difference between the top and base voltage
// of specified source
func (d *DSOX3024) GetVoltageAmplitude(channel int) (float64, error) {
	response, err := d.Query(fmt.Sprintf(":MEASure:VAMPlitude? CHANnel%d", channel))
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(response), 64)
}

// Calculates the average voltage over the displayed waveform
func (d *DSOX3024) GetVoltageAverage(channel int) (float64, error) {
	response, err := d.Query(fmt.Sprintf(":MEASure:VAVerage? DISPlay, CHANnel%d", channel))
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(response), 64)
}

// Converts a byte array to an array of float64 values (SWAP, 64-bit)
func (d *DSOX3024) ByteToFloatArray(payload []byte) ([]float64, error) {
	if len(payload)%8 != 0 {
		return nil, fmt.Errorf("payload must be aligned to 8 bytes, got %d", len(payload))
	}
	numValues := len(payload) / 8

	values := make([]float64, numValues)
	for i := 0; i < numValues; i++ {
		bits := binary.LittleEndian.Uint64(payload[i*8 : i*8+8])
		values[i] = math.Float64frombits(bits)
	}

	return values, nil
}

// Sets the waveform points mode (RAW, NORMal, MAXimum)
func (d *DSOX3024) SetWaveformPointsMode(mode string) error {
	return d.Write(fmt.Sprintf(":WAVeform:POINts:MODE %s", mode))
}

// Sets the waveform number of points
func (d *DSOX3024) SetWaveformPoints(points int) error {
	return d.Write(fmt.Sprintf(":WAVeform:POINts %d", points))
}

// Get the waveform from a desired channel
func (d *DSOX3024) GetWaveform(channel int) (*Waveform, error) {
	if err := d.WriteSequence([]string{
		fmt.Sprintf(":WAV:SOUR CHANnel%d", channel),
		":WAVeform:FORMat BYTE",
		":WAVeform:UNSigned OFF",
	}); err != nil {
		return nil, err
	}

	data, err := d.QueryByteSequence(":WAV:DATA?")
	if err != nil {
		return nil, err
	}

	p, err := d.GetPreamble(channel)
	if err != nil {
		return nil, err
	}
	values := make([]float64, len(data))
	for i, b := range data {
		adc := float64(b)
		values[i] = (adc-p.yReference)*p.yIncrement + p.yOrigin
	}
	timestamps := make([]float64, len(values))
	for i := range timestamps {
		timestamps[i] = p.xOrigin + (float64(i)-p.xReference)*p.xIncrement
	}

	return &Waveform{
		Values:     values,
		Timestamps: timestamps,
	}, nil
}
