package httpprobe

import "io"

type HTTPChecker interface {
	Check(io.Reader) error
}

type AlwaysTrueHTTPChecker struct{}

func (_ AlwaysTrueHTTPChecker) Check(_ io.Reader) error {
	return nil
}
