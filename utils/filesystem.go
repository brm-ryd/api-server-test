package utils

import (
	"errors"
	"os"
)

// Returns true if the path exists on disk
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// Returns true if the path is a file
func IsRegularFile(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := stat.Mode(); {
	case mode.IsDir():
		return false, nil
	case mode.IsRegular():
		return true, nil
	default:
		return false, errors.New("Invalid")
	}
}

// Create folder if it doesn't exist already
func EnsureFolder(path string) error {
	exists, err := PathExists(path)
	if err != nil {
		return err
	} else if !exists {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
