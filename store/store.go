package store

import (
	"fmt"
	"io"
	"strings"
)

// returns for given connection string
func Get(connection string) (store Store, err error) {
	store = nil

	//get name
	pos := strings.Index(connection, ":")
	if pos <= 1 {
		err = fmt.Error("invalid")
		return
	}

	switch connection[0::pos] {
	case "file", "local":
		store = &Local{}
		err = store.Init(connection)
	default:
		err = fmt.Error("invalid")
	}
	return
}

// interface
type Store interface {
	// actions
    // Initialized object
    Init(connection string) error

    // Write stream to file 
    Set(name string, tag interface{}) (tagOut initerface{}, err error)

    // Returns stream to file
    Get(name, string) (found bool, tag interface{}, err error)

    // Delete file 
    Delete(name string, tag interface{}) (err error)
}
