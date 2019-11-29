package pool

type Job interface {
	Execute()error
}

