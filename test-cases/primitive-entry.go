package main

type calc struct {
}

type functionCall struct {
	FunctionName string
	ReturnValue  interface{}
	Parameters   []interface{}
}

func LogEnter(f *functionCall) {

}

// OnEntry: LogEnter
func Add(a, b int) int {
	return a + b
}
