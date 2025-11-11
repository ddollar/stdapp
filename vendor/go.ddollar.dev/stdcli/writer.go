package stdcli

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"go.ddollar.dev/errors"
)

var (
	DefaultWriter *Writer
)

type Renderer func(string) string

type Writer struct {
	Color  bool
	Stdout io.Writer
	Stderr io.Writer
	Tags   map[string]Renderer
}

func init() {
	DefaultWriter = &Writer{
		Color:  isTerminal(os.Stdout),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tags: map[string]Renderer{
			"error":  renderError,
			"header": RenderColors(242),
			"h1":     RenderColors(244),
			"h2":     RenderColors(241),
			"id":     RenderColors(247),
			"info":   RenderColors(247),
			"ok":     RenderColors(46),
			"start":  RenderColors(247),
			"u":      RenderUnderline(),
			"value":  RenderColors(251),
		},
	}
}

func (w *Writer) Error(err error) error {
	fmt.Fprintf(w.Stderr, w.renderTags("<error>%s</error>\n"), err)

	if os.Getenv("DEBUG") == "true" {
		if serr, ok := err.(errors.ErrorTracer); ok {
			for _, f := range serr.ErrorTrace() {
				fmt.Fprintf(w.Stderr, w.renderTags("<info>  %s:%d</info>\n"), f, f)
			}
		}
	}

	return err //nowrap
}

func (w *Writer) Errorf(format string, args ...any) error {
	return w.Error(errors.Errorf(format, args...))
}

func (w *Writer) IsTerminal() bool {
	if f, ok := w.Stdout.(*os.File); ok {
		return isTerminal(f)
	}

	return false
}

func (w *Writer) Sprintf(format string, args ...any) string {
	return fmt.Sprintf(w.renderTags(format), args...)
}

func (w *Writer) Write(data []byte) (int, error) {
	n, err := w.Stdout.Write([]byte(w.renderTags(string(data))))
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return n, nil
}

func (w *Writer) Writef(format string, args ...any) (int, error) {
	n, err := fmt.Fprintf(w.Stdout, w.renderTags(format), args...)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	return n, nil
}

func (w *Writer) renderTags(s string) string {
	for tag, render := range w.Tags {
		s = regexp.MustCompile(fmt.Sprintf("<%s>(.*?)</%s>", tag, tag)).ReplaceAllStringFunc(s, render)
	}

	if !w.Color {
		s = stripColor(s)
	}

	return s
}

func RenderColors(colors ...int) Renderer {
	return func(s string) string {
		s = stripTag(s)
		for _, c := range colors {
			s = fmt.Sprintf("\033[38;5;%dm", c) + s
		}
		return s + "\033[0m"
	}
}

func RenderUnderline() Renderer {
	return func(s string) string {
		return fmt.Sprintf("\033[4m%s\033[24m", stripTag(s))
	}
}

func renderError(s string) string {
	return fmt.Sprintf("\033[38;5;124mERROR: \033[38;5;203m%s\033[0m", stripTag(s))
}

var (
	colorStripper = regexp.MustCompile("\033\\[[^m]+m")
	tagStripper   = regexp.MustCompile(`^<[^>?]+>(.*)</[^>?]+>$`)
	tagMatcher    = regexp.MustCompile(`<([^>?]+)>`)
)

func stripColor(s string) string {
	return colorStripper.ReplaceAllString(s, "")
}

func stripTag(v any) string {
	s := fmt.Sprintf("%v", v)

	match := tagStripper.FindStringSubmatch(s)

	if len(match) != 2 {
		return s
	}

	return match[1]
}

func stripTags(v any) string {
	s := fmt.Sprintf("%v", v)

	for {
		m := tagMatcher.FindStringSubmatchIndex(s)

		if len(m) != 4 {
			break
		}

		os := m[0]
		oe := m[1]

		closer := fmt.Sprintf("</%s>", s[m[2]:m[3]])
		cs := strings.Index(s, closer)

		if cs == -1 {
			break
		}

		ce := cs + len(closer)

		s = s[:os] + s[oe:cs] + s[ce:]
	}

	return s
}
