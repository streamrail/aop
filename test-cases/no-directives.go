package main

type calc struct {
}

func Add(a, b int) int {
	return a + b
}

func (c calc) Sub(a, b int) int {
	return a - b
}

func (c *calc) Mul(a, b int) int {
	return a * b
}

func Div(a, b []*int) int {
	return a / b
}
