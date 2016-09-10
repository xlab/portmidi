package main

import (
	"math"
	"sync"
)

type DSP interface {
	Process(out []float32)
}

type sineOsc struct {
	step  float64
	phase float64
}

func NewSineOsc(freq float64, sampleRate float64) DSP {
	return &sineOsc{freq / sampleRate, 0}
}

func (osc *sineOsc) Process(out []float32) {
	for i := range out {
		out[i] = osc.powerInNthSample(i)
	}
}

func (osc *sineOsc) powerInNthSample(n int) float32 {
	currentPhase := osc.phase
	_, osc.phase = math.Modf(osc.phase + osc.step)
	return float32(math.Sin(2 * math.Pi * currentPhase))
}

type Vocoder struct {
	note int
	dsps map[int]DSP
	mux  sync.Mutex
}

func NewVocoder() *Vocoder {
	v := &Vocoder{
		dsps: make(map[int]DSP),
	}
	v.SwitchNote(60) // C4
	return v
}

func noteToFreq(n int) float64 {
	// http://newt.phys.unsw.edu.au/jw/notes.html
	return math.Pow(2, (float64(n)-69)/12) * 440.0
}

func (v *Vocoder) SwitchNote(note int) {
	v.mux.Lock()
	v.note = note
	_, ok := v.dsps[note]
	if !ok {
		v.dsps[note] = NewSineOsc(noteToFreq(note), 44100)
	}
	v.mux.Unlock()
}

func (v *Vocoder) CurrentDSP() DSP {
	v.mux.Lock()
	dsp := v.dsps[v.note]
	v.mux.Unlock()
	return dsp
}
