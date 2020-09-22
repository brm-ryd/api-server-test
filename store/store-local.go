package store

import (
	"io"
	"os"
	"path"
	// "os/user"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/brm-ryd/api-server-test/utils"
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

	path, err := os.UserHomeDir(path)
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
