package stdapp

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

//go:generate zip -r init.zip init

//go:embed init.zip
var initArchive []byte

func initApp(name string) error {
	if err := os.Mkdir(name, 0755); err != nil {
		return errors.Wrap(err, "mkdir")
	}

	if err := os.Chdir(name); err != nil {
		return errors.Wrap(err, "chdir")
	}

	if err := exec.Command("git", "init").Run(); err != nil {
		return errors.Wrap(err, "git init")
	}

	zr, err := zip.NewReader(bytes.NewReader(initArchive), int64(len(initArchive)))
	if err != nil {
		return err
	}

	app, err := fs.Sub(zr, "init")
	if err != nil {
		return errors.Wrap(err, "chroot")
	}

	err = fs.WalkDir(app, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		target := path

		if strings.HasSuffix(path, ".tmpl") {
			target = strings.TrimSuffix(target, ".tmpl")
		}

		fmt.Println(target)

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		data, err := fs.ReadFile(app, path)
		if err != nil {
			return err
		}

		data = bytes.ReplaceAll(data, []byte("stdapp-init"), []byte(name))
		data = bytes.ReplaceAll(data, []byte("example.org/stdapp"), []byte(fmt.Sprintf("github.com/ddollar/%s", name)))

		if err := ioutil.WriteFile(target, data, 0644); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "walk")
	}

	if err := os.Symlink("docker-compose.development.yml", "docker-compose.override.yml"); err != nil {
		return err
	}

	if err := exec.Command("make", "vendor").Run(); err != nil {
		return errors.Wrap(err, "vendor")
	}

	return nil
}
