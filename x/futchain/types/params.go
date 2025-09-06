package types

const DefaultTimezone string = "Europe/Istanbul"

// NewParams creates a new Params instance.
func NewParams(timezone string) Params {
	return Params{Timezone: timezone}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultTimezone)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := validateTimezone(p.Timezone); err != nil {
		return err
	}

	return nil
}
func validateTimezone(v string) error {

	return nil
}
