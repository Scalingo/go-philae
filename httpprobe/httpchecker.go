package httpprobe

import "io"

type HTTPChecker interface {
	Check(io.Reader) error
}
