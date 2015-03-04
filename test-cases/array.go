package main

import (
	"fmt"
	"time"
)

type calc struct {
}

type functionCall struct {
	FunctionName string
	ReturnValue  interface{}
	Parameters   []interface{}
}

var cache map[string]interface{}

func LogEnter(f *functionCall) {
	fmt.Printf("LogEnter %s\n", f.FunctionName)
}

func LogReturn(f *functionCall) {
	fmt.Printf("LogReturn %s\n", f.FunctionName)
}

func GetFromCache(f *functionCall) {
	if val, ok := cache[f.FunctionName]; ok {
		f.ReturnValue = val
	}
}

func StoreInCache(f *functionCall) {
	cache[f.FunctionName] = f.ReturnValue
}

// OnEntry: LogEnter
// OnEntry: GetFromCache
// OnReturn: LogReturn
// OnReturn: StoreInCache
func Add(a, b []int) []int {
	return []int{a[0] + b[0]}
}

func main() {
	cache = make(map[string]interface{})
	a := []int{1}
	b := []int{2}

	fmt.Printf("START %v\n", time.Now())
	fmt.Printf("FIRST_CALL %v;%d\n", time.Now(), Add(a, b)[0])
	fmt.Printf("SECOND_CALL %v;%d\n", time.Now(), Add(a, b)[0])
}
