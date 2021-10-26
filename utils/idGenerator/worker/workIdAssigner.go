package worker

type WorkerIdAssigner interface {
	AssignWorkerId() int64
	GetNamespace() string
}