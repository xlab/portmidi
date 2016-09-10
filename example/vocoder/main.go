package main

import (
	"flag"
	"log"

	"github.com/xlab/closer"
	"github.com/xlab/midievent"
	"github.com/xlab/portmidi"
)

var (
	inName  = flag.String("in", "", "MIDI device name to use as input.")
	inDevID = flag.Int("in-dev", -1, "MIDI device ID to use as input.")
)

func init() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
}

const samplesPerChannel = 2048

func main() {
	defer closer.Close()

	portmidi.Initialize()
	closer.Bind(func() {
		portmidi.Terminate()
	})

	numDevices := portmidi.CountDevices()
	if numDevices < 1 {
		closer.Fatalln("[ERR] vocoder cannot operate with less than one MIDI device")
	}
	inputs := make(map[string]portmidi.DeviceID, numDevices)
	for i := 0; i < numDevices; i++ {
		info := portmidi.GetDeviceInfo(portmidi.DeviceID(i))
		if info.IsInputAvailable {
			inputs[info.Name] = portmidi.DeviceID(i)
		}
	}
	log.Println("[INFO] available inputs:", inputs)
	devID := findCandidate(inputs, *inDevID, *inName, true)
	info := portmidi.GetDeviceInfo(devID)
	log.Printf("Using %s (via %s)", info.Name, info.Interface)
	midiIn, err := portmidi.NewInputStream(devID, 512, 0,
		portmidi.FilterControl|portmidi.FilterAftertouch|portmidi.FilterSystemCommon|portmidi.FilterRealtime)
	if err != nil {
		closer.Fatalln(err)
	}
	closer.Bind(func() {
		midiIn.Close()
	})

	vocoder := NewVocoder()
	go func() {
		for ev := range midiIn.Source() {
			msg := portmidi.Message(ev.Message)
			if midievent.IsNoteOn(midievent.Event(msg.Status())) {
				n := int(msg.Data1())
				log.Printf("note %d (%.3fHz)", n, noteToFreq(n))
				vocoder.SwitchNote(n)
			}
		}
	}()

	in := make(chan []float32, 64)
	out := make(chan []float32, 64)
	go func() {
		buf := make([]float32, 2*samplesPerChannel)

		for frame := range in {
			dsp := vocoder.CurrentDSP()
			dsp.Process(buf)
			for i := range frame {
				buf[i] = (frame[i] * buf[i])
			}
			out <- buf
			buf = make([]float32, 2*samplesPerChannel)
		}
	}()
	ctl, err := NewIOControl(2, 2, samplesPerChannel, in, out)
	if err != nil {
		closer.Fatalln(err)
	}
	closer.Bind(func() {
		ctl.Destroy()
	})

	if err := ctl.StartStream(); err != nil {
		closer.Fatalln(err)
	}
	closer.Hold()
}

func findCandidate(devices map[string]portmidi.DeviceID,
	id int, name string, input bool) (dev portmidi.DeviceID) {

	if id >= 0 {
		dev = portmidi.DeviceID(id)
		return
	}
	if len(name) > 0 {
		nameID, ok := devices[name]
		if ok {
			dev = nameID
			return
		}
		closer.Fatalln("[ERR] vocoder was unable to locate required device:", name)
	}
	if input {
		dev, _ = portmidi.DefaultInputDeviceID()
		return
	}
	dev, _ = portmidi.DefaultOutputDeviceID()
	return
}
