package runtime

type Validatable interface {
	Validate() error // Validate validates the struct
}
