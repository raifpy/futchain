package types

const DefaultTimezone string = "Europe/Istanbul"
const DefaultFetchModulo int64 = 3

// NewParams creates a new Params instance.
func NewParams(timezone string, fetchModulo int64) Params {
	return Params{Timezone: timezone, FetchModulo: fetchModulo}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultTimezone, DefaultFetchModulo)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := validateTimezone(p.Timezone); err != nil {
		return err
	}
	if err := validateFetchModulo(p.FetchModulo); err != nil {
		return err
	}

	return nil
}
func validateTimezone(v string) error {

	return nil
}
func validateFetchModulo(v int64) error {

	return nil
}
