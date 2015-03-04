package main

type functionCall struct {
	FunctionName string
	ReturnValue  interface{}
	Parameters   []interface{}
}

func LogEnter(f *functionCall) {

}

func LogReturn(f *functionCall) {

}

// OnEntry: LogEnter
// OnReturn: LogReturn
func Add(a int) *int {
	return &a
}
