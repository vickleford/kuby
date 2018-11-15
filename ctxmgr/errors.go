package ctxmgr

type ErrConfigParse struct {
	what string
}

func (e *ErrConfigParse) Error() string {
	return e.what
}

func newErrConfigParse(what string) *ErrConfigParse {
	return &ErrConfigParse{what: what}
}
