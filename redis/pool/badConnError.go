package pool

type BadConnError struct {
	wrapped error
}

var _ error = (*BadConnError)(nil)

func (e BadConnError) Error() string {
	s := "redis: Conn is in a bad state"
	if e.wrapped != nil {
		s += ": " + e.wrapped.Error()
	}
	return s
}

func (e BadConnError) Unwrap() error {
	return e.wrapped
}
