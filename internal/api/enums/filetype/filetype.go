package filetype

// FileType is the type of a file
type FileType int

const (
	Pdf FileType = iota
	Word
	Txt
)

func (f FileType) String() string {
	return [...]string{"pdf", "word", "txt"}[f]
}
