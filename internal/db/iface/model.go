package iface

import "github.com/driif/echo-go-starter/internal/server/net/runtime"

type Model interface {
	FromDTO(dto runtime.Validatable)
}
