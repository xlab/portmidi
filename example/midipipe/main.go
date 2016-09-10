package main

import (
	"flag"
	"log"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/xlab/closer"
	"github.com/xlab/portmidi"
	"github.com/xlab/treeprint"
)

var (
	inName   = flag.String("in", "", "MIDI device name to use as input.")
	outName  = flag.String("out", "", "MIDI device name to use as output.")
	inDevID  = flag.Int("in-dev", -1, "MIDI device ID to use as input.")
	outDevID = flag.Int("out-dev", -1, "MIDI device ID to use as output.")
)

func init() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)
}

func main() {
	defer closer.Close()

	portmidi.Initialize()
	closer.Bind(func() {
		portmidi.Terminate()
	})

	numDevices := portmidi.CountDevices()
	log.Println("[INFO] total MIDI devices:", numDevices)
	if numDevices < 2 {
		closer.Fatalln("[ERR] midipipe cannot operate with less than 2 devices")
	}
	inputs := make(map[string]portmidi.DeviceID, numDevices)
	outputs := make(map[string]portmidi.DeviceID, numDevices)
	for i := 0; i < numDevices; i++ {
		info := portmidi.GetDeviceInfo(portmidi.DeviceID(i))
		if info.IsInputAvailable {
			inputs[info.Name] = portmidi.DeviceID(i)
		}
		if info.IsOutputAvailable {
			outputs[info.Name] = portmidi.DeviceID(i)
		}
	}
	log.Println("[INFO] available inputs:", inputs)
	log.Println("[INFO] available outputs:", outputs)
	inDev := findCandidate(inputs, *inDevID, *inName, true)
	outDev := findCandidate(outputs, *outDevID, *outName, false)

	inInfo := portmidi.GetDeviceInfo(inDev)
	log.Printf("[INFO] input device id=%d %s", inDev, treeprint.Repr(inInfo))
	in, err := portmidi.NewInputStream(inDev, 1024, 0)
	if err != nil {
		closer.Fatalln("[ERR] cannot init an input stream:", err)
	}
	closer.Bind(func() {
		in.Close()
	})
	outInfo := portmidi.GetDeviceInfo(outDev)
	log.Printf("[INFO] output device id=%d %s", outDev, treeprint.Repr(outInfo))
	out, err := portmidi.NewOutputStream(outDev, 1024, 0, 0)
	if err != nil {
		closer.Fatalln("[ERR] cannot init an output stream:", err)
	}
	closer.Bind(func() {
		out.Close()
		log.Println("bye!")
	})

	meter := metrics.NewMeter()
	go func() {
		t := time.NewTicker(time.Minute)
		for range t.C {
			snap := meter.Snapshot()
			log.Println("[DBG] rate 1 minute", snap.Rate1())
			log.Println("[DBG] rate 5 minute", snap.Rate5())
			log.Println("[DBG] rate mean", snap.RateMean())
		}
	}()
	go func() {
		sink := out.Sink()
		for ev := range in.Source() {
			meter.Mark(int64(1))
			sink <- ev
		}
	}()
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
		closer.Fatalln("[ERR] midipipe was unable to locate required device:", name)
	}
	if input {
		dev, _ = portmidi.DefaultInputDeviceID()
		return
	}
	dev, _ = portmidi.DefaultOutputDeviceID()
	return
}
