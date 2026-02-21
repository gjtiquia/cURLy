package input

type Buffer []Action

type Action byte

const (
	ActionNone Action = iota
	ActionUp
	ActionDown
	ActionLeft
	ActionRight
	ActionRestart
	ActionExit
)

func CreateBuffer() Buffer {
	// assuming player wont make more than 8 inputs per frame
	return make(Buffer, 0, 8)
}
