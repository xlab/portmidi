portmidi
========

The package provides Go bindings for [PortMIDI](http://portmedia.sourceforge.net/portmidi/) from the PortMedia set of libraries. [PortMedia](http://portaudio.com) offers free, cross-platform, open-source I/O libraries for digital media including MIDI, video, and audio.<br/>
All the binding code for the `pm` package has automatically been generated with rules defined in [pm.yml](/pm.yml). The wrapping package `portmidi` has been done by hand and leverages channels for MIDI event streaming.

## Usage

```
$ brew install portmidi
(or use your package manager)

$ go get github.com/xlab/portmidi
```

## Examples

### MIDIPipe

```
$ brew install portmidi
$ go get github.com/xlab/portmidi/example/midipipe

$ midipipe -in "Arturia BeatStep" -out "OP-1 Midi Device"

main.go:35: [INFO] total MIDI devices: 4
main.go:50: [INFO] available inputs: map[Arturia BeatStep:1 OP-1 Midi Device:0]
main.go:51: [INFO] available outputs: map[OP-1 Midi Device:2 Arturia BeatStep:3]

main.go:56: [INFO] input device id=1 .
├── [CoreMIDI]  Interface
├── [Arturia BeatStep]  Name
├── [true]  IsInputAvailable
└── [false]  IsOutputAvailable
main.go:65: [INFO] output device id=2 .
├── [CoreMIDI]  Interface
├── [OP-1 Midi Device]  Name
├── [false]  IsInputAvailable
└── [true]  IsOutputAvailable

main.go:80: [DBG] rate 1 minute 100.14564014611834
main.go:81: [DBG] rate 5 minute 51.93229225191669
main.go:82: [DBG] rate mean 134.09919032344553

main.go:80: [DBG] rate 1 minute 65.02414338183416
main.go:81: [DBG] rate 5 minute 51.98104984601001
main.go:82: [DBG] rate mean 93.71848260518783
^Cmain.go:72: bye!
```

`midipipe` is simple Go program that redirects all the events it gets from the MIDI input device
to the specified MIDI output device. You can specify route by device name (see example) or by its ID.

The app requires minimum two devices to operate properly, but note that a single hardware piece can act
both as input and output device, so by "devices" I mean logical I/O streams.

### Vocoder

```
$ brew install portaudio portmidi
$ go get github.com/xlab/portmidi/example/vocoder

$ vocoder -in "Arturia BeatStep"
main.go:43: [INFO] available inputs: map[IAC Driver Bus 1:0 Arturia BeatStep:1]
main.go:46: Using Arturia BeatStep (via CoreMIDI)
main.go:62: note 63 (311.127Hz)
main.go:62: note 62 (293.665Hz)
main.go:62: note 65 (349.228Hz)
main.go:62: note 73 (554.365Hz)
main.go:62: note 66 (369.994Hz)
main.go:62: note 60 (261.626Hz)
^C
```

`vocoder` is an implementation of a simple vocoder in Go. It reads your voice using PortAudio, reads note-on events from your MIDI device using PortMIDI, and plays the altered voice back using PortAudio. Have fun.

### Rebuilding the package

You will need to get the [cgogen](https://git.io/cgogen) tool installed first.

```
$ git clone https://github.com/xlab/portmidi && cd portmidi
$ make clean
$ make
```

## License

All the code except when stated otherwise is licensed under the MIT license.
