package testutils

import (
	"testing"
)

type Wrapper struct {
	*testing.T
}

func (w Wrapper) BailIfError(err error) {
	if err != nil {
		w.Fatal(err)
	}
}

func (w Wrapper) BailIfErrorf(err error, msg string, args ...interface{}) {
	if err != nil {
		w.Fatalf(msg, args...)
	}
}

func (w Wrapper) ShouldFailf(err error, msg string, args ...interface{}) {
	if err == nil {
		w.Fatalf(msg, args...)
	}
}

func Wrap(t *testing.T) Wrapper {
	return Wrapper{t}
}
