package util

type foo string

func (f foo) Clone() foo {
	return f
}
