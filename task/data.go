package task

type Type string

type Data interface {
	ID() uint64
	Type() Type
	Data() []byte
}
