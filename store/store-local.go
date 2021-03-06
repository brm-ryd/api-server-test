package store

import (
	"io"
	"os"
	"path"
	// "os/user"
	"errors"
	"fmt"
	"path/filepath"
	// "runtime"
	"strings"

	"github.com/brm-ryd/api-server-test/utils"
	homedir "github.com/mitchellh/go-homedir"
)

// local file system
type Local struct {
	basePath string
}

func (f *Local) Init(connection string) error {
	// Ensure that connection starts with "local:" or "file:"
	if !strings.HasPrefix(connection, "local:") && !strings.HasPrefix(connection, "file:") {
		return fmt.Errorf("invalid scheme")
	}

	path := connection[strings.Index(connection, ":")+1:]

	path, err := homedir.Expand(path)
	if err != nil {
		return err
	}

	// Get the path
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}

	// Ensure the path ends with a /
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// Lastly, ensure the path exists
	err = utils.EnsureFolder(path)
	if err != nil {
		return err
	}

	f.basePath = path

	return nil
}

func (f *Local) Get(name string, out io.Writer) (found bool, tag interface{}, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}

	found = true

	// Opening file
	file, err := os.Open(f.basePath + name)
	if err != nil {
		if os.IsNotExist(err) {
			found = false
			err = nil
		}
		return
	}
	defer file.Close()

	// Check file content
	stat, err := file.Stat()
	if err != nil {
		return
	}
	if stat.Size() == 0 {
		found = false
		return
	}

	// Copy file to stream
	_, err = io.Copy(out, file)
	if err != nil {
		return
	}

	return
}

func (f *Local) Set(name string, in io.Reader, tag interface{}) (tagOut interface{}, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}

	// Create folders if necessary
	dir := path.Dir(name)
	if dir != "" {
		err = os.MkdirAll(f.basePath+dir, os.ModePerm)
		if err != nil {
			return
		}
	}

	// Create file
	file, err := os.Create(f.basePath + name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Write stream to file
	_, err = io.Copy(file, in)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (f *Local) Delete(name string, tag interface{}) (err error) {
	if name == "" {
		err = errors.New("empty")
		return
	}

	// Delete the file
	err = os.Remove(f.basePath + name)
	return
}
