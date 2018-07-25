package embed

import (
	"github.com/gobuffalo/packr"
	"github.com/juju/errors"
)

// registry handles bindata
var registry = make(map[string]packr.Box)

// Register initializes bindata and stores it in the registry
func Register(group string) {
	switch group {
	case "migrations":
		Make("app/migrations/default", packr.NewBox("./../../app/migrations/default"))
		return
	}

	return
}

// Make stores the bindata in the registry
func Make(name string, box packr.Box) (err error) {
	if Exists(name) {
		err = errors.Errorf("Embed already exists: %", name)
		return
	}

	registry[name] = box
	return nil
}

// Get fetches the bindata from the registry
func Get(name string) (box packr.Box, err error) {
	if !Exists(name) {
		err = errors.Errorf("Embed does not exist: %", name)
		return
	}

	box = registry[name]
	return
}

// Exists checks if the bindata is stored in the registry
func Exists(name string) bool {
	if _, ok := registry[name]; ok {
		return true
	}

	return false
}
