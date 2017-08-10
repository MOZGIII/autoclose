package autoclose

import "io"

// Reader reads data
type Reader struct {
	upstream   io.ReadCloser
	onlyOnEOF  bool
	onCloseErr func(closeErr, readErr error) error
}

// NewReader provides a simple way to build an autoclose reader.
func NewReader(upstream io.ReadCloser) *Reader {
	return &Reader{upstream, true, ReturnReadErr}
}

// ReturnReadErr simply returns readErr.
// This is used with an autoclose reader as onCloseErr.
func ReturnReadErr(closeErr, readErr error) error {
	return readErr
}

func (r Reader) Read(p []byte) (n int, err error) {
	n, err = r.upstream.Read(p)
	if err != nil {
		if !r.onlyOnEOF || err == io.EOF {
			if closeErr := r.upstream.Close(); closeErr != nil {
				err = r.onCloseErr(closeErr, err)
			}
		}
	}
	return
}

// Compile-time interface check
var _ io.Reader = &Reader{}
