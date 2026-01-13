package export

import "io"

type Exporter[T any] interface {
	ContentType() string
	FileName() string
	Export(w io.Writer, data []T) error
}
