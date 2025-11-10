package stdcli

import (
	"io"
	"os"

	"golang.org/x/term"
)

var (
	DefaultReader *Reader
)

type Reader struct {
	io.Reader
}

func init() {
	DefaultReader = &Reader{os.Stdin}
}

func (r *Reader) IsTerminal() bool {
	if f, ok := r.Reader.(*os.File); ok {
		return isTerminal(f)
	}

	return false
}

func (r *Reader) TerminalRaw() func() bool {
	var fd int
	var state *term.State

	if f, ok := r.Reader.(*os.File); ok {
		fd = int(f.Fd())
		if s, err := term.MakeRaw(fd); err == nil {
			state = s
		}
	}

	return func() bool {
		if state != nil {
			if err := term.Restore(fd, state); err != nil {
				return false
			}

			return true
		}
		return false
	}
}
