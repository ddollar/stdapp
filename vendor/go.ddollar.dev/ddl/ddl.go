package ddl

import (
	"io"

	"golang.org/x/exp/constraints"
)

func If[T any](condition bool, thenValue, elseValue T) T {
	if condition {
		return thenValue
	}

	return elseValue
}

func Mean[T constraints.Float | constraints.Integer](values []T) T {
	if len(values) == 0 {
		return T(0)
	}

	var sum T
	for _, v := range values {
		sum += v
	}

	return sum / T(len(values))
}

type multiWriteCloser struct {
	io.Writer
	cs []io.Closer
}

func MultiWriteCloser(ws ...io.Writer) io.WriteCloser {
	m := &multiWriteCloser{Writer: io.MultiWriter(ws...)}

	for _, w := range ws {
		if c, ok := w.(io.Closer); ok {
			m.cs = append(m.cs, c)
		}
	}

	return m
}

func (m *multiWriteCloser) Close() error {
	var first error

	for _, c := range m.cs {
		if err := c.Close(); err != nil && first == nil {
			first = err
		}
	}

	return first
}
