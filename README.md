# aop
AOP process GO source files looking for specific function decorated comments which directs it to wrap such functions with calls to OnEntry and OnReturn


##example

Suppose that src.go contains the following Add function

```go

	// OnEntry: LogEnter
	// OnEntry: GetFromCache
	// OnReturn: LogReturn
	// OnReturn: StoreInCache
	func Add(a, b int) int {
		return a + b
	}

```

The Add function is decorated with OnEntry and OnReturn, each decoration specifies
which function(s) it wishes to call, for the example above LogEnter and GetFromCache will be called immediately (in this order) each time Add is called, once the Add function returns both LogReturn and StoreInCache are called (in the order specified).

A decorating function has the following signature:

```go
	func (f *functionCall)
```

Where functionCall is a struct:

```go
	type functionCall struct {
	FunctionName string
	ReturnValue  interface{}
	Parameters   []interface{}
}
```

functionName: holds the name of the function which causes this decoration to run, in our case Add.

ReturnValue: is nil for OnEntry 'event', within OnReturn 'event' it will hold the return value of the callee

Parameters: holds the parameters which the callee has been called with.

You may alter both Parameters and ReturnValue, please note setting ReturnValue to a none nil value will cause the callee (in this example Add) to return that value immediately without executing.


Calling
```cmd
	aop src.go output.go
```

Will generate a new output.go file:

```go
func OriginalAdd(a, b int) int {
	return a + b
}

func Add(a, b) {
	ctx := &functionCall{
		FunctionName: "Add",
		ReturnValue:  nil,
		Parameters:   []interface{}{a, b},
	}

	LogEnter(ctx)
	GetFromCache(ctx)

	if ctx.ReturnValue != nil {
		LogReturn(ctx)
		StoreInCache(ctx)

		return ctx.ReturnValue.(int)
	}

	ctx.ReturnValue = OriginalAdd(a, b)

	LogReturn(ctx)
	StoreInCache(ctx)

	return ctx.ReturnValue.(int)
}

```

output.go renamed 'Add' to 'OriginalAdd' and created an entire new 'Add' function

Please note AOP is in its early stages, it can deal with a subset of function signatures, we're working on expending its supports for wider range of function signatures.
