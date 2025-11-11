package stdapi

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	templates       FileSystem
	templateHelpers TemplateHelpers
)

// TemplateHelpers is a function that returns template helper functions for a given request.
//
// The Context is provided so helpers can access request-specific data like sessions or user info.
type TemplateHelpers func(c *Context) template.FuncMap

// FileSystem is an alias for http.FileSystem used for serving static files and templates.
type FileSystem http.FileSystem

// LoadTemplates configures the template system with a filesystem and optional helper functions.
//
// The FileSystem is typically an http.Dir() pointing to your templates directory.
// TemplateHelpers can be nil if no custom functions are needed.
//
// Example:
//
//	stdapi.LoadTemplates(http.Dir("./templates"), func(c *stdapi.Context) template.FuncMap {
//		return template.FuncMap{
//			"formatDate": func(t time.Time) string { return t.Format("2006-01-02") },
//		}
//	})
func LoadTemplates(files FileSystem, helpers TemplateHelpers) {
	templates = files
	templateHelpers = helpers
}

// TemplateExists checks if a template file exists in the configured filesystem.
func TemplateExists(path string) bool {
	_, err := templates.Open(path)
	return !os.IsNotExist(err)
}

// RenderTemplate renders an HTML template with hierarchical layout resolution.
//
// This function automatically searches for layout.tmpl files in parent directories,
// allowing templates to inherit from layouts. The search proceeds from the root
// to the template's directory.
//
// For example, rendering "admin/users/list" will load templates in this order:
//  1. layout.tmpl (root)
//  2. admin/layout.tmpl
//  3. admin/users/layout.tmpl
//  4. admin/users/list.tmpl
//
// Templates should define a "main" block that layouts can invoke.
func RenderTemplate(c *Context, path string, params interface{}) error {
	return RenderTemplatePart(c, path, "main", params)
}

// RenderTemplatePart renders a specific named block from a template.
//
// The part parameter specifies which template block to execute (default is "main").
func RenderTemplatePart(c *Context, path, part string, params interface{}) error {
	files := []string{}

	files = append(files, "layout.tmpl")

	parts := strings.Split(filepath.Dir(path), "/")

	for i := range parts {
		files = append(files, filepath.Join(filepath.Join(parts[0:i+1]...), "layout.tmpl"))
	}

	files = append(files, fmt.Sprintf("%s.tmpl", path))

	ts := template.New(part)

	if templateHelpers != nil {
		ts = ts.Funcs(templateHelpers(c))
	}

	for _, f := range files {
		fd, err := templates.Open(f)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}

		if _, err := ts.Parse(string(data)); err != nil {
			return errors.WithStack(err)
		}
	}

	var buf bytes.Buffer

	if err := ts.Execute(&buf, params); err != nil {
		return errors.WithStack(err)
	}

	io.Copy(c, &buf)

	return nil
}
