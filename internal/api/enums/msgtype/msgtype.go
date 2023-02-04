package msgtype

// MsgType is the type of a message
type MsgType int

const (
	Text MsgType = iota
	Image
	Video
	Audio
	File
)

func (m MsgType) String() string {
	return [...]string{"text", "image", "video", "audio", "file"}[m]
}
