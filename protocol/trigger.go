package protocol

import "fmt"

// Sets the trigger mode (AUTO, NORMal)
func (d *DSOX3024) SetTriggerMode(mode string) error {
	cmd := fmt.Sprintf(":TRIGger:MODE %s", mode)
	return d.Write(cmd)
}
