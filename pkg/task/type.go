package task

type Type byte

const (
	TypeMap Type = iota
	TypeReduce
	TypeNone
)
