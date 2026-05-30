package trash

import "errors"

var (
	ErrNoBranchOrMatchProvided = errors.New("no branch or match provided")
	ErrInvalidDayCount         = errors.New("expected 1 or more days provided")
)
