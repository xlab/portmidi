package pm

import (
	"errors"
	"fmt"
)

var True Error = 1
var False Error = 0

const NoDevice = nodevice

var (
	ErrHostError = errors.New("portmidi: host error")
	// ErrInvalidDeviceID means out of range or
	// output device when input is requested or
	// input device when output is requested or
	// device is already opened.
	ErrInvalidDeviceID    = errors.New("portmidi: invalid DeviceID")
	ErrInsufficientMemory = errors.New("portmidi: insufficient memory")
	ErrBufferTooSmall     = errors.New("portmidi: buffer too small")
	ErrBufferOverflow     = errors.New("portmidi: buffer overflow")
	// ErrBadPtr is a result of PortMidiStream parameter being nil or
	// stream is not opened or
	// stream is output when input is required or
	// stream is input when output is required.
	ErrBadPtr = errors.New("portmidi: bad pointer")
	// ErrBadData means illegal midi data, e.g. missing EOX.
	ErrBadData       = errors.New("portmidi: bad data")
	ErrInternalError = errors.New("portmidi: internal error")
	// ErrBufferMaxSize means buffer is already as large as it can be.
	ErrBufferMaxSize = errors.New("portmidi: buffer max size")
	ErrUnknown       = errors.New("portmidi: unknown error")
)

func HasData(e Error) bool {
	if e == gotdata {
		return true
	}
	return false
}

func ToError(e Error) error {
	switch e {
	case noerror, gotdata:
		return nil
	case hosterror:
		return ErrHostError
	case invaliddeviceid:
		return ErrInvalidDeviceID
	case insufficientmemory:
		return ErrInsufficientMemory
	case buffertoosmall:
		return ErrBufferTooSmall
	case bufferoverflow:
		return ErrBufferOverflow
	case badptr:
		return ErrBadPtr
	case baddata:
		return ErrBadData
	case internalerror:
		return ErrInternalError
	case buffermaxsize:
		return ErrBufferMaxSize
	default:
		return fmt.Errorf("portmidi: %s", GetErrorText(e))
	}
}
