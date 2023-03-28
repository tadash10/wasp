package verify

import "fmt"

type ErrInvalidURLScheme struct {
	ValidSchemes  []string
	InvalidScheme string
}

func NewErrInvalidURLScheme(invalidScheme string) error {
	return ErrInvalidURLScheme{
		ValidSchemes:  validSchemes,
		InvalidScheme: invalidScheme,
	}
}

func (e ErrInvalidURLScheme) Error() string {
	return fmt.Sprintf("%s is not a valid URL scheme, must be one of %v", e.InvalidScheme, e.ValidSchemes)
}

func (e ErrInvalidURLScheme) Is(target error) bool {
	_, ok := target.(ErrInvalidURLScheme)
	return ok
}

type ErrInvalidURLHostname struct {
	InvalidHostname string
}

func NewErrInvalidURLHostname(invalidHostname string) error {
	return ErrInvalidURLHostname{
		InvalidHostname: invalidHostname,
	}
}

func (e ErrInvalidURLHostname) Error() string {
	if len(e.InvalidHostname) == 0 {
		return "URL must include a hostname"
	}
	return fmt.Sprintf("%s is not a valid hostname", e.InvalidHostname)
}

func (e ErrInvalidURLHostname) Is(target error) bool {
	_, ok := target.(ErrInvalidURLHostname)
	return ok
}
