package pool

import "sync"

type CondStatus uint32

const (
	RUN CondStatus = iota
	PAUSE
	EXIT
)

type Cond struct {
	*sync.Cond
	status CondStatus
}

func NewCond() *Cond {
	return &Cond{
		Cond:   sync.NewCond(&sync.Mutex{}),
		status: PAUSE,
	}
}

func NewCondWithLock(l sync.Locker) *Cond {
	return &Cond{
		Cond:   sync.NewCond(l),
		status: PAUSE,
	}
}

func NewCondWithStatus(l sync.Locker, status CondStatus) *Cond {
	return &Cond{
		Cond:   sync.NewCond(l),
		status: status,
	}
}
