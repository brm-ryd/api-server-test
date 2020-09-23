package store

import (
	"fmt"
	"io"
	"strings"
)

// returns for given connection string
func Get(connection string) (store Store, err error) {
	//get name
	pos := strings.Index(connection, ":")

}

// interface
type Store interface {
	// actions
}
