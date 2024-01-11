package wrr

import "errors"

var (
	ErrUncompleted   = errors.New("uncompleted")
	ErrIllegalParams = errors.New("illegal params")
	ErrEmptyData     = errors.New("empty data")
)
