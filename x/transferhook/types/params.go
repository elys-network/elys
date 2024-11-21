package types

const (
	// Default active value for each transferhook supported module
	DefaultAmmActive = true
)

// NewParams creates a new Params instance
func NewParams(ammActive bool) Params {
	return Params{
		AmmActive: ammActive,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(DefaultAmmActive)
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}
