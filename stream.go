package portmidi

import (
	"time"

	"github.com/xlab/portmidi/pm"
)

type Stream struct {
	stream *pm.PortMidiStream
	buf    chan Event
	closeC chan struct{}
	doneC  chan struct{}
}

// Close closes a midi stream, flushing any pending buffers.
func (s *Stream) Close() error {
	close(s.closeC)
	<-s.doneC
	err := pm.ToError(pm.Close(s.stream))
	s.stream = nil
	return err
}

func (s *Stream) Source() <-chan Event {
	return s.buf
}

func (s *Stream) Sink() chan<- Event {
	return s.buf
}

// NewInputStream opens device for the input. The buffersize specifies the number of input events to be
// buffered waiting to be read.
func NewInputStream(id DeviceID, bufferSize int,
	channels ChannelMask, filters ...Filter) (*Stream, error) {

	var stream *pm.PortMidiStream
	ret := pm.OpenInput(&stream, pm.DeviceID(id), nil, int32(bufferSize), nil, nil)
	if err := pm.ToError(ret); err != nil {
		return nil, err
	}
	buf := make(chan Event, bufferSize)
	s := &Stream{
		stream: stream,
		closeC: make(chan struct{}),
		doneC:  make(chan struct{}),
		buf:    buf,
	}
	if channels > 0 { // all allowed by default
		pm.SetChannelMask(s.stream, int32(channels))
	}
	if len(filters) > 0 {
		var fs Filter
		fs.Join(filters...)
		pm.SetFilter(s.stream, int32(fs))
	}
	go s.processInput()
	return s, nil
}

// NewOutputStream opens device for the input. The buffersize
// specifies the number of output events to be buffered waiting for output.
// (In some cases -- see below -- PortMidi does not buffer output at all
// and merely passes data to a lower-level API, in which case buffersize
// is ignored.)
//
// latency is the delay in milliseconds applied to timestamps to determine
// when the output should actually occur. (If latency is < 0, 0 is assumed.)
// If latency is zero, timestamps are ignored and all output is delivered
// immediately. If latency is greater than zero, output is delayed until the
// message timestamp plus the latency. (NOTE: the time is measured relative
// to the time source indicated by time_proc. Timestamps are absolute,
// not relative delays or offsets.) In some cases, PortMidi can obtain
// better timing than your application by passing timestamps along to the
// device driver or hardware. Latency may also help you to synchronize midi
// data to audio data by matching midi latency to the audio buffer latency.
func NewOutputStream(id DeviceID, bufferSize, latency int,
	channels ChannelMask, filters ...Filter) (*Stream, error) {
	var stream *pm.PortMidiStream
	ret := pm.OpenOutput(&stream, pm.DeviceID(id), nil, int32(bufferSize), nil, nil, int32(latency))
	if err := pm.ToError(ret); err != nil {
		return nil, err
	}
	buf := make(chan Event, bufferSize)
	s := &Stream{
		stream: stream,
		closeC: make(chan struct{}),
		doneC:  make(chan struct{}),
		buf:    buf,
	}
	if channels > 0 { // all allowed by default
		pm.SetChannelMask(s.stream, int32(channels))
	}
	go s.processOutput()
	return s, nil
}

func (s *Stream) pushEvents(buf []pm.Event) {
	for i := range buf {
		if buf[i].Ref() == nil {
			return
		}
		buf[i].Deref()
		s.buf <- Event{
			Timestamp: int32(buf[i].Timestamp),
			Message:   Message(buf[i].Message),
		}
	}
}

const pollDelay = 5 * time.Millisecond

func (s *Stream) processInput() {
	var hadData bool
	for {
		select {
		case <-s.closeC:
			close(s.buf)
			close(s.doneC)
			return
		default:
			if pm.Poll(s.stream) == pm.True {
				hadData = true
				buf := make([]pm.Event, 0, 512)
				size := pm.Read(s.stream, buf, 512)
				s.pushEvents(buf[:size])
				continue
			} else if hadData {
				hadData = false
				continue
			}
			time.Sleep(pollDelay)
		}
	}
}

// An aggregating version of this function available:
// https://gist.github.com/xlab/1768c3dd210bf3829b54f4cec3f748bb
func (s *Stream) processOutput() {
	for {
		select {
		case <-s.closeC:
			go func() {
				// drain s.buf
				for range s.buf {
				}
			}()
			close(s.doneC)
			return
		case ev, ok := <-s.buf:
			if !ok { // s.buf closed
				close(s.doneC)
				return
			}
			if len(ev.SysExData) > 0 { // handle sysEx separately
				pm.WriteSysEx(s.stream, pm.Timestamp(ev.Timestamp), ev.SysExData)
				continue
			}
			pm.WriteShort(s.stream, pm.Timestamp(ev.Timestamp), int32(ev.Message))
		}
	}
}

// HasHostError tests whether stream has a pending host error.
// Normally, the client finds out about errors through returned error codes,
// but some errors can occur asynchronously where the client does not
// explicitly call a function, and therefore cannot receive an error code.
func (s *Stream) HasHostError() bool {
	return pm.HasHostError(s.stream) > 0
}

// Synchronize instructs PortMidi to (re)synchronize to the
// time_proc passed when the stream was opened.
// PortMidi will always synchronize at the
// first output message and periodically thereafter.
// func (s *Stream) Sync() error {
// 	return pm.ToError(pm.Synchronize(s.stream))
// }
