package diagrams

import (
	"github.com/emarais-godaddy/go-diagrams/diagram"
)

func New(opts ...diagram.Option) (*diagram.Diagram, error) {
	return diagram.New(opts...)
}
