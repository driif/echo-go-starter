package videotype

// VideoType is the type of a video
type VideoType int

const (
	Mp4 VideoType = iota
	Webm
)

func (v VideoType) String() string {
	return [...]string{"mp4", "webm"}[v]
}
