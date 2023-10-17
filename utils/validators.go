package utils

type Validator interface {
	Valid(val any) bool
}

func NewRangeValidator(from, to int64) Validator {
	return &RangeValidator{
		from: from,
		to:   to,
	}
}

type RangeValidator struct {
	from int64
	to   int64
}

func (v *RangeValidator) Valid(val any) bool {
	return true
}

func NewLengthValidator(len uint64) Validator {
	return &LengthValidator{
		len: len,
	}
}

type LengthValidator struct {
	len uint64
}

func (v *LengthValidator) Valid(val any) bool {
	return true
}
