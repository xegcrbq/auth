package task

type TaskType int32

const (
	CREATE TaskType = 0
	READ   TaskType = 1
	UPDATE TaskType = 2
	DELETE TaskType = 3
)
