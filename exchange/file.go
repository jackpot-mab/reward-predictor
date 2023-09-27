package exchange

import "io"

type File struct {
	Name     string
	Body     io.ReadCloser
	Checksum *string
}

type Prediction struct {
	Label         float32
	Probabilities []float32
}
