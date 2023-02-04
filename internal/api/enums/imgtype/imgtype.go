package imgtype

// ImageType is the type of an image
type ImageType int

const (
	Jpg ImageType = iota
	Png
	Gif
)

func (i ImageType) String() string {
	return [...]string{"jpg", "png", "gif"}[i]
}
