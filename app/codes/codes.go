package codes

import (
	"net/http"

	"github.com/epointpayment/mloc-cpe/app/log"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Code contains information about an application code
type Code struct {
	Name       string
	Message    string
	Group      string
	StatusCode int
	IsError    bool
}

// New creates an application code
func New(name string) *Code {
	return &Code{
		Name: name,
	}
}

// WithMessage attaches a message for the code
func (ac *Code) WithMessage(message string) *Code {
	ac.Message = message
	return ac
}

// WithStatusCode sets the HTTP status code
func (ac *Code) WithStatusCode(statusCode int) *Code {
	ac.StatusCode = statusCode
	return ac
}

// InGroup adds group information
func (ac *Code) InGroup(name string) *Code {
	ac.Group = name
	return ac
}

// RegisterMessage adds an application code message entry to the registry
func (ac *Code) RegisterMessage() Code {
	if ac.StatusCode == 0 {
		ac.StatusCode = http.StatusOK
	}

	return ac.register()
}

// RegisterError adds an application code error entry to the registry
func (ac *Code) RegisterError() Code {
	if ac.StatusCode == 0 {
		ac.StatusCode = http.StatusInternalServerError
	}

	return ac.register()
}

// Register adds an application code entry to the registry
func (ac *Code) register() Code {
	res := *ac

	if err := res.Validate(); err != nil {
		log.Fatal(res.Name+" is an invalid application code.\n", err)
	}

	if err := registry.Add(res); err != nil {
		log.Fatal(err)
	}

	return res
}

// Validate checks the application code for invalid values
func (ac Code) Validate() (err error) {
	err = validation.ValidateStruct(&ac,
		validation.Field(&ac.Name, validation.Required),
		validation.Field(&ac.Message, validation.Required),
		validation.Field(&ac.StatusCode, validation.Required),
	)

	return
}

// Error returns a string describing the error
func (ac Code) Error() string {
	if !ac.IsError {
		return ""
	}

	return ac.Message
}

// Get gets an application code from the application code registry
func Get(code string) (ac Code, err error) {
	ac, err = registry.Get(code)
	if err != nil {
		return
	}

	return
}
