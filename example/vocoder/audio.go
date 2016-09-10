package main

import (
	"errors"
	"unsafe"

	"github.com/xlab/portaudio-go/portaudio"
)

type IOControl struct {
	stream *portaudio.Stream
	inCh   int
	outCh  int
	in     chan<- []float32
	out    <-chan []float32
}

const (
	paFormat     = portaudio.PaFloat32
	paSampleRate = 44100
)

func NewIOControl(inCh, outCh int, samplesPerCh int,
	in chan<- []float32, out <-chan []float32) (*IOControl, error) {
	if err := portaudio.Initialize(); paError(err) {
		return nil, errors.New(paErrorText(err))
	}
	if in == nil {
		inCh = 0
	}
	if out == nil {
		outCh = 0
	}
	var stream *portaudio.Stream
	ctl := &IOControl{
		inCh:  inCh,
		outCh: outCh,
		in:    in,
		out:   out,
	}
	if err := portaudio.OpenDefaultStream(&stream, int32(inCh), int32(outCh), paFormat, paSampleRate,
		uint(samplesPerCh), ctl.audioCallback, nil); paError(err) {
		return nil, errors.New(paErrorText(err))
	}
	ctl.stream = stream
	return ctl, nil
}

func (i *IOControl) StartStream() error {
	err := portaudio.StartStream(i.stream)
	if paError(err) {
		return errors.New(paErrorText(err))
	}
	return nil
}

func (i *IOControl) audioCallback(input unsafe.Pointer, output unsafe.Pointer, sampleCount uint,
	_ *portaudio.StreamCallbackTimeInfo, _ portaudio.StreamCallbackFlags, _ unsafe.Pointer) int32 {

	const statusContinue = int32(portaudio.PaContinue)
	samples := int(sampleCount)

	if input != nil && i.in != nil {
		inFrame := (*(*[1 << 32]float32)(input))[:samples*i.inCh]
		i.in <- inFrame[:samples*i.inCh] // TODO(xlab): consider copying
	}
	if output != nil && i.out != nil {
		outFrame := (*(*[1 << 32]float32)(output))[:samples*i.outCh]
		select {
		case frame := <-i.out:
			copy(outFrame, frame[:samples*i.outCh])
		default:
			return statusContinue
		}
	}
	return statusContinue
}

func (i *IOControl) Destroy() {
	if i.stream != nil {
		portaudio.StopStream(i.stream)
		portaudio.Terminate()
		i.stream = nil
	}
}

func paError(err portaudio.Error) bool {
	return portaudio.ErrorCode(err) != portaudio.PaNoError
}

func paErrorText(err portaudio.Error) string {
	return portaudio.GetErrorText(err)
}
