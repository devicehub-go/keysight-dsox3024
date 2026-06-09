package dsox3024

import (
	"github.com/devicehub-go/keysight-dsox3024/protocol"
	"github.com/devicehub-go/unicomm"
)

/*
Creates a new instance of Keysight E5061B Network
Analyzer to communicate and control the device
*/
func New(options unicomm.Options) *protocol.DSOX3024 {
	options.Delimiter = "\n"
	return &protocol.DSOX3024{
		Communication: unicomm.New(options),
	}
}
