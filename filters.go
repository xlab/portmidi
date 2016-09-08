package portmidi

type ChannelMask int32

// Channel sets a channel mask to pass to a stream.
// Multiple channels should be OR'ed together Channel(1) | Channel(2).
// Note that channels are numbered 0 to 15 (not 1 to 16).
// All channels are allowed by default.
func Channel(ch int) ChannelMask {
	return ChannelMask(1 << uint8(ch))
}

type Filter int32

func (f *Filter) Join(fs ...Filter) {
	for i := range fs {
		*f |= fs[i]
	}
}

const (
	// FilterActive filters active sensing messages (0xFE).
	FilterActive Filter = (1 << 0x0E)
	// FilterSysEx filters system exclusive messages (0xF0).
	FilterSysEx Filter = (1 << 0x00)
	// FilterClock filters MIDI clock messages (0xF8).
	FilterClock Filter = (1 << 0x08)
	// FilterPlay filters play messages (start 0xFA, stop 0xFC, continue 0xFB).
	FilterPlay Filter = ((1 << 0x0A) | (1 << 0x0C) | (1 << 0x0B))
	// FilterTick filters tick messages.
	FilterTick Filter = (1 << 0x09)
	// FilterFd filters undefined FD messages.
	FilterFD Filter = (1 << 0x0D)
	// FilterUndefined filters undefined real-time messages.
	FilterUndefined Filter = FilterFD
	// FilterReset filters reset messages.
	FilterReset Filter = (1 << 0x0F)
	// FilterRealtime filters all real-time messages.
	FilterRealtime Filter = (FilterActive | FilterSysEx | FilterClock | FilterPlay | FilterUndefined | FilterReset | FilterTick)
	// FilterNote filters note-on and note-off (0x90-0x9F and 0x80-0x8F).
	FilterNote Filter = ((1 << 0x19) | (1 << 0x18))
	// FilterChannelAftertouch filters channel aftertouch (most midi controllers use this) (0xD0-0xDF).
	FilterChannelAftertouch Filter = (1 << 0x1D)
	// FilterPolyAftertouch filters per-note aftertouch (0xA0-0xAF).
	FilterPolyAftertouch Filter = (1 << 0x1A)
	// FilterAftertouch filters both channel and poly aftertouch.
	FilterAftertouch Filter = (FilterChannelAftertouch | FilterPolyAftertouch)
	// FilterProgram filters program changes (0xC0-0xCF).
	FilterProgram Filter = (1 << 0x1C)
	// FilterControl filters control changes (CC's) (0xB0-0xBF).
	FilterControl Filter = (1 << 0x1B)
	// FilterPitchbend filters pitch bends (0xE0-0xEF).
	FilterPitchbend Filter = (1 << 0x1E)
	// FilterMTC filters MIDI Time Code (0xF1).
	FilterMTC Filter = (1 << 0x01)
	// FilterSongPosition filters song position (0xF2).
	FilterSongPosition Filter = (1 << 0x02)
	// FilterSongSelect filters song select (0xF3).
	FilterSongSelect Filter = (1 << 0x03)
	// FilterTune filters tuning request (0xF6).
	FilterTune Filter = (1 << 0x06)
	// FilterSystemCommon filters all system common messages (mtc, song position, song select, tune request).
	FilterSystemCommon Filter = (FilterMTC | FilterSongPosition | FilterSongSelect | FilterTune)
)
