package protocol

import (
	"fmt"
	"slices"
	"strconv"
)

// Sets the vertical scale, in volts units per division,
// of the selected channel
func (d *DSOX3024) SetVerticalScale(channel int, units float64) error {
	cmd := fmt.Sprintf(":CHANnel%d:SCALe %f", channel, units)
	return d.Write(cmd)
}

// Gets the vertical scale, in volts units per division
func (d *DSOX3024) GetVerticalScale(channel int) (float64, error) {
	response, err := d.Query(fmt.Sprintf(":CHANnel%d:SCALe?", channel))
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(response), 64)
}

// Sets the time interval between the trigger event
// and the delay reference point
func (d *DSOX3024) SetTimebasePosition(position float64) error {
	cmd := fmt.Sprintf(":TIMebase:POSition %f", position)
	return d.Write(cmd)
}

// Gets the current delay value in seconds
func (d *DSOX3024) GetTimebasePosition() (float64, error) {
	response, err := d.Query(":TIMebase:POSition?")
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(response), 64)
}

// Sets the full-scale horizontal time in seconds
func (d *DSOX3024) SetTimebaseRange(fullScale float64) error {
	cmd := fmt.Sprintf(":TIMebase:RANGe %f", fullScale)
	return d.Write(cmd)
}

// Gets the full-scale horizontal time in seconds
func (d *DSOX3024) GetTimebaseRange() (float64, error) {
	response, err := d.Query(":TIMebase:RANGe?")
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(response), 64)
}

// Sets the horizontal reference position to a percent-of-screen
// location, from left to right
func (d *DSOX3024) SetTimebaseReference(location string) error {
	options := []string{"LEFT", "RIGHt", "CENTer"}
	if !slices.Contains(options, location) {
		return fmt.Errorf("reference must LEFT, RIGHt or CENTer, got %s", location)
	}
	cmd := fmt.Sprintf(":TIMebase:REFerence %s", location)
	return d.Write(cmd)
}

// Gets the horizontal reference position in percentage-of-screen
func (d *DSOX3024) GetTimebaseReference() (string, error) {
	response, err := d.Query(":TIMebase:REFerence?")
	return string(response), err
}

// Sets the time base scale in seconds per divison
func (d *DSOX3024) SetTimebaseScale(units float64) error {
	cmd := fmt.Sprintf(":TIMebase:SCALe %f", units)
	return d.Write(cmd)
}
