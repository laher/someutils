package someutils

type Signal interface {
	Status() int
}

type SignalSimple struct {
	StatusCode int
}

func (ss *SignalSimple) Status() int {
	return ss.StatusCode
}

var (
	SIGINT = &SignalSimple{1}
)
