package codes

import "github.com/juju/errors"

// registry contains a list of application codes
var registry = NewRegistry()

// Registry manages the applications codes
type Registry struct {
	store map[string]Code
}

// NewRegistry initializes the registry for use
func NewRegistry() *Registry {
	return &Registry{
		store: make(map[string]Code),
	}
}

// Get gets an application code entry
func (r *Registry) Get(code string) (ac Code, err error) {
	if v, ok := r.store[code]; ok {
		ac = v
		return
	}

	err = errors.New("code does not exist in registry: " + code)
	return
}

// Exists checks if an application code entry exists in the list
func (r *Registry) Exists(code string) bool {
	if _, ok := r.store[code]; ok {
		return true
	}

	return false
}

// Add adds an application code entry to the list
func (r *Registry) Add(ac Code) error {
	code := ac.Name
	if !r.Exists(code) {
		r.store[code] = ac
		return nil
	}

	return errors.New("code already exists in registry: " + ac.Name)
}

// Remove removes an application code entry to the list
func (r *Registry) Remove(code string) bool {
	if r.Exists(code) {
		delete(r.store, code)
		return true
	}

	return false
}
