package pool

import (
	"fmt"
	"github.com/astaxie/beego"
)

type BasicTask struct {
	Name string
	Pool *WorkPool
}

func (t *BasicTask) QueuedWork() int32 {
	return t.Pool.QueuedWork()
}

func (t *BasicTask) PreDoWork(workRoutine int) {
	qw := t.Pool.QueuedWork()
	ar := t.Pool.ActiveRoutines()
	beego.Debug(fmt.Sprintf("*******> Task: %s WR: %d QW: %d AR: %d Total: %d\n",
		t.Name,
		workRoutine,
		qw,
		ar,
		qw+ar))
}
