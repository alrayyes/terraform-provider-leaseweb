package enum_utils

// Used in FindEnumForInt to get the int value of an enum.
type intEnum interface {
	// Value returns the enum int value.
	Value() int
}
