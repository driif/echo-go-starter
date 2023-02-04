package audiotype

// AudioType is the type of an audio
type AudioType int

const (
	Mp3 AudioType = iota
	Wav
)

func (a AudioType) String() string {
	return [...]string{"mp3", "wav"}[a]
}
