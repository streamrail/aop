package main

type calc struct {
}

type functionCall struct {
	FunctionName string
	ReturnValue  interface{}
	Parameters   []interface{}
}

func LogReturn(f *functionCall) {

}

// OnReturn: LogReturn
func Add(a, b int) int {
	return a + b
}
