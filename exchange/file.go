package exchange

import "io"

type File struct {
	Name     string
	Body     io.ReadCloser
	Checksum *string
}

type Prediction struct {
	Label         int
	Probabilities []float32
}
