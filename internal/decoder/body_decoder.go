package decoder

import (
	"encoding/json"
	"io"
)

func NewDecoder[T any](r io.Reader) *Decoder[T] {
	return &Decoder[T]{r: r}
}

type Decoder[T any] struct {
	r io.Reader
}

func (d *Decoder[T]) Decode(v *T) error {
	return json.NewDecoder(d.r).Decode(&v)
}

func (d *Decoder[T]) Result() (v T, err error) {
	err = d.Decode(&v)
	return
}
