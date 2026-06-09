package protocol

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/devicehub-go/unicomm"
)

type DSOX3024 struct {
	Communication unicomm.Unicomm
	mutex         sync.Mutex
}

// Establishes a connection with the device
func (d *DSOX3024) Connect() error {
	if err := d.Communication.Connect(); err != nil {
		return err
	} else if err := d.Communication.Write([]byte("*IDN?\n")); err != nil {
		return err
	}
	response, err := d.Communication.ReadUntil("\n")
	if err != nil {
		return err
	}
	fmt.Println("Connected to ", string(response))
	return nil
}

// Closes the connection with the device
func (d *DSOX3024) Disconnect() error {
	return d.Communication.Disconnect()
}

// Returns true if device is connected
func (d *DSOX3024) IsConnected() bool {
	return d.Communication.IsConnected()
}

// Writes a message to device
func (d *DSOX3024) Write(message string) error {
	if !d.IsConnected() {
		return fmt.Errorf("no device connected")
	}
	if !strings.Contains(message, "\n") {
		message = message + "\n"
	}
	return d.Communication.Write([]byte(message))
}

// Writes a sequence of messages to device
func (d *DSOX3024) WriteSequence(messages []string) error {
	for _, message := range messages {
		if err := d.Write(message); err != nil {
			return err
		}
	}
	return nil
}

// Reads byte response from device
func (d *DSOX3024) Query(message string) ([]byte, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.Write(message); err != nil {
		return nil, err
	}
	response, err := d.Communication.ReadUntil("\n")
	if err != nil {
		return nil, err
	}
	return response[:len(response)-1], nil
}

// Reads a sequence of bytes from device
func (d *DSOX3024) QueryByteSequence(message string) ([]byte, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.Write(message); err != nil {
		return nil, err
	}
	header, err := d.Communication.Read(2)
	if err != nil {
		return nil, err
	}
	if len(header) != 2 || header[0] != '#' {
		return nil, fmt.Errorf("invalid header, got %s", string(header))
	}
	numDigits := int(header[1] - '0')
	numBytes, err := d.Communication.Read(uint(numDigits))
	if err != nil {
		return nil, err
	}
	n, err := strconv.Atoi(string(numBytes))
	if err != nil {
		return nil, err
	}
	d.Communication.Read(1)

	payload := make([]byte, 0)
	for len(payload) != int(n) {
		response, err := d.Communication.Read(uint(n - len(payload)))
		if err != nil {
			return nil, err
		}
		payload = append(payload, response...)
	}
	return payload, nil
}
