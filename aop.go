package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"strings"
)

// TODO: not sure if we should use LEAVE(ctx) in "if ctx.ReturnValue != nil { "
const (
	cTemplate = `SIG {
	ctx := &functionCall{
		FunctionName: "FNAME",
		ReturnValue:  nil,
		Parameters: []interface{} {PARAMS},
	}

	ENTRY(ctx)
	if ctx.ReturnValue != nil { 
		LEAVE(ctx)
		return ctx.ReturnValue.(RET_TYPE)
	}

	ctx.ReturnValue = ORIGINAL(PARAMS)

	LEAVE(ctx)

	return ctx.ReturnValue.(RET_TYPE)
}
`
)

// Discribes a function parameter
type Parameter struct {
	Name string // Parameter's name
	Type string // Parameter's type
}

func NewParameter() *Parameter {
	return &Parameter{}
}

var functions map[string]*FuncDesc // Holds decorated functions.

// Groups function declaration with its directive.
type FuncDesc struct {
	F *ast.FuncDecl
	D []*Directives
}

func NewFuncDesc() *FuncDesc {
	return &FuncDesc{
		F: nil,
		D: make([]*Directives, 0),
	}
}

// Function's directives.
type Directives struct {
	OnEntry  string // Function to call upon entering the decorated function (// OnEntry: function_name)
	OnReturn string // Function to call upon returning from the decorated function (// OnReturn: function_name)
}

// Looks for directives.
// comment: a single line comment.
func MatchDirectives(comment string) *Directives {
	idx := strings.Index(comment, "OnEntry: ")
	if idx != -1 {
		idx += len("OnEntry: ")
		return &Directives{OnEntry: comment[idx:]}

	}

	idx = strings.Index(comment, "OnReturn: ")
	if idx != -1 {
		idx += len("OnReturn: ")
		return &Directives{OnReturn: comment[idx:]}
	}

	return nil
}

type FuncVisitor struct {
}

func NewFuncVisitor() *FuncVisitor {
	return &FuncVisitor{}
}

// Visitor to traverse an abstract syntax tree.
func (v *FuncVisitor) Visit(node ast.Node) (w ast.Visitor) {

	// which type are we dealing with?
	switch t := node.(type) {
	case *ast.FuncDecl:
		if t.Doc != nil {
			if len(t.Doc.List) > 0 {

				f := NewFuncDesc()

				// Scan through function's comments
				for _, comment := range t.Doc.List {

					// See if comment is a directive
					if match := MatchDirectives(comment.Text); match != nil {
						f.D = append(f.D, match)
					}
				}

				// Have we found any directives?
				if len(f.D) > 0 {
					f.F = t
					// Store function.
					functions[t.Name.Name] = f
					// Rename function.
					t.Name.Name = "Original" + t.Name.Name
				}

				return nil // No need to visit sub nodes.
			}
		}
	}

	return v // Keep going.
}

// Returns function arguments as a comma delimited string.
func GetParams(ft *ast.FuncType) string {
	params := ""

	// Extract input parameters.
	for _, param := range ft.Params.List {
		for _, name := range param.Names {
			params += name.Name + ", "
		}
	}

	// Remove last ', '
	params = params[:len(params)-2]
	return params
}

// Returns function arguments.
func GetParamsWithTypes(ft *ast.FuncType) []*Parameter {
	return ParseParamList(ft.Params.List)
}

// Returns function's return arguments
// TODO this function needs to be recursive, as type can be nested *[]*[] etc'
func GetReturnTypes(ft *ast.FuncType) []*Parameter {

	/*Results: *ast.FieldList {
	  236  .  .  .  .  .  Opening: 0
	  237  .  .  .  .  .  List: []*ast.Field (len = 1) {
	  238  .  .  .  .  .  .  0: *ast.Field {
	  239  .  .  .  .  .  .  .  Type: *ast.StarExpr {
	  240  .  .  .  .  .  .  .  .  Star: 255
	  241  .  .  .  .  .  .  .  .  X: *ast.ArrayType {
	  242  .  .  .  .  .  .  .  .  .  Lbrack: 256
	  243  .  .  .  .  .  .  .  .  .  Elt: *ast.StarExpr {
	  244  .  .  .  .  .  .  .  .  .  .  Star: 258
	  245  .  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
	  246  .  .  .  .  .  .  .  .  .  .  .  NamePos: 259
	  247  .  .  .  .  .  .  .  .  .  .  .  Name: "int"
	  248  .  .  .  .  .  .  .  .  .  .  }
	  249  .  .  .  .  .  .  .  .  .  }
	  250  .  .  .  .  .  .  .  .  }
	  251  .  .  .  .  .  .  .  }
	  252  .  .  .  .  .  .  }
	  253  .  .  .  .  .  }
	  254  .  .  .  .  .  Closing: 0
	  255  .  .  .  .  }*/
	params := make([]*Parameter, 0)

	// Extract return types.
	for _, field := range ft.Results.List {
		param := NewParameter()
		param.Name = ""

		switch t := field.Type.(type) {
		case *ast.ArrayType:
			switch eltType := t.Elt.(type) {
			case *ast.Ident:
				param.Type = "[]" + eltType.Name
			}
		case *ast.Ident:
			param.Type = t.Name

		case *ast.StarExpr:
			switch xType := t.X.(type) {
			case *ast.Ident:
				param.Type = "*" + xType.Name
			}

		default:
			fmt.Printf("Unhandled type: %s\n", t)
		}
		params = append(params, param)
	}

	return params
}

func ParseParamList(parameters []*ast.Field) []*Parameter {
	params := make([]*Parameter, 0)

	// Extract input parameters.
	for _, param := range parameters {
		ps := make([]*Parameter, 0)
		for _, name := range param.Names {
			p := NewParameter()
			p.Name = name.Name
			ps = append(ps, p)
		}

		// Determin type.
		pType := ""
		switch tt := param.Type.(type) {
		case *ast.ArrayType:
			switch eltType := tt.Elt.(type) {
			case *ast.Ident:
				pType = "[]" + eltType.Name
			}
		case *ast.Ident:
			pType = tt.Name
		}

		for _, p := range ps {
			p.Type = pType
		}

		params = append(params, ps...)
	}

	return params
}

// Constructs function signature.
func BuildSignature(reciver *Parameter, name string, arguments, returnParams []*Parameter) string {
	args := ""
	rets := ""

	if len(arguments) > 0 {
		for _, a := range arguments {
			args += a.Name + " " + a.Type + ", "
		}
		args = args[:len(args)-2]
	}

	if len(returnParams) > 0 {
		for _, r := range returnParams {
			rets += r.Type + ", "
		}
		rets = rets[:len(rets)-2]
	}

	strReciver := ""
	if reciver != nil {
		strReciver = fmt.Sprintf("(%s %s)", reciver.Name, reciver.Type)
	}

	if len(returnParams) < 2 {
		return fmt.Sprintf("func %s %s(%s) %s", strReciver, name, args, rets)
	} else {
		return fmt.Sprintf("func %s %s(%s) (%s)", strReciver, name, args, rets)
	}
}

// Incase function is a Method, reciver is the type on which the function belongs to.
func GetReciver(recv *ast.FieldList) *Parameter {
	p := NewParameter()
	if recv != nil {
		if recv.List != nil {
			field := recv.List[0]
			p.Name = field.Names[0].Name
			switch fieldType := field.Type.(type) {
			case *ast.StarExpr:
				switch fieldTypeXType := fieldType.X.(type) {
				case *ast.Ident:
					p.Type = "*" + fieldTypeXType.Name
				}
			case *ast.Ident:
				p.Type = fieldType.Name
			}
			return p
		}
	}
	return nil
}

// Creates a new wrapper function which wraps decorated function,
// This new wrapper function is responsible to call each directive (OnEntry, OnReturn).
// func WrapFunc(funcName string, fd *FuncDesc) *ast.FuncDecl {
func WrapFunc(funcName string, fd *FuncDesc) string {
	fn := fd.F
	template := string(cTemplate)

	returnParams := GetReturnTypes(fn.Type)
	arguments := GetParamsWithTypes(fn.Type)
	reciver := GetReciver(fn.Recv)
	sig := BuildSignature(reciver, funcName, arguments, returnParams)

	if len(returnParams) > 1 {
		fmt.Println("Currently we don't support functions which return multiple argumetns, this will be supported in the near future")
		return ""
	}

	// Replace.
	template = strings.Replace(template, "SIG", sig, 1)
	template = strings.Replace(template, "FNAME", funcName, 1)

	entries := ""
	returns := ""
	for _, directive := range fd.D {
		if len(directive.OnEntry) > 0 {
			entries += directive.OnEntry + "(ctx)\n"
		}

		if len(directive.OnReturn) > 0 {
			returns += directive.OnReturn + "(ctx)\n"
		}
	}

	if len(entries) > 0 {
		template = strings.Replace(template, "ENTRY(ctx)", entries, 1)
	} else {
		template = strings.Replace(template, "ENTRY(ctx)", "", 1)
		template = strings.Replace(template, "	if ctx.ReturnValue != nil { return ctx.ReturnValue.(RET_TYPE) }", "", 1)
	}

	if len(returns) > 0 {
		template = strings.Replace(template, "LEAVE(ctx)", returns, -1)
	} else {
		template = strings.Replace(template, "LEAVE(ctx)", "", -1)
		template = strings.Replace(template, "ctx.ReturnValue = ORIGINAL(PARAMS)", "return ORIGINAL(PARAMS)", 1)
		template = strings.Replace(template, "	return ctx.ReturnValue.(RET_TYPE)", "", 1)
	}

	if reciver != nil && reciver.Name != "" {
		template = strings.Replace(template, "ORIGINAL", reciver.Name+"."+fn.Name.Name, 1) // fn.Name.Name already replace. (prefixed with Original)
	} else {
		template = strings.Replace(template, "ORIGINAL", fn.Name.Name, -1) // fn.Name.Name already replace. (prefixed with Original)
	}

	template = strings.Replace(template, "PARAMS", GetParams(fn.Type), -1)
	template = strings.Replace(template, "RET_TYPE", returnParams[0].Type, -1)

	return template

	// err := ioutil.WriteFile("instance.go", []byte(template), 0644)

	// Parse.
	// fset := token.NewFileSet()
	// f, err := parser.ParseFile(fset, "instance.go", nil, parser.ParseComments)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }

	// ast.Print(nil, f)

	// TODO, instade of function declerations to the abs simply append
	// the result of printer Fprint to the bottom of our output file.
	// printer.Fprint(os.Stdout, fset, f)
	// return f.Decls[0].(*ast.FuncDecl)
}

func ProcessFile(src, output string) error {
	functions = make(map[string]*FuncDesc)
	fset := token.NewFileSet()

	// Parse source file
	f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Prints ABS tree.
	// ast.Print(nil, f)

	// Locate decorated functions
	ast.Walk(NewFuncVisitor(), f)

	outputFile, err := os.Create(output)
	if err != nil {
		fmt.Printf("error creating %s\n", output)
		return err
	}
	defer outputFile.Close()

	// clones src to dest.
	printer.Fprint(outputFile, fset, f)

	wrappers := ""
	// Wrap each function
	for funcName, desc := range functions {
		wrapper := WrapFunc(funcName, desc)

		// Add wrapper function to ast.
		// f.Decls = append(f.Decls, wrapper)
		wrappers += "\n" + wrapper
	}

	// Append wrapper function to the end of the output file.
	if err = AppendToFile(output, wrappers); err != nil {
		fmt.Printf("error appending wrapper functions to %s %v\n", output, err)
		return err
	}

	if err = FormatSource(output); err != nil {
		fmt.Printf("error formating source %s %v\n", output, err)
		return err
	}

	return nil
	// Prints ABS tree.
	// ast.Print(nil, f)
	// printer.Fprint(os.Stdout, fset, f)
}

func Usage() {
	fmt.Println("-src source.go -out /out/source.go")
}

func main() {
	src := flag.String("src", "", "source file to process")
	out := flag.String("out", "", "output file")
	flag.Parse()

	fmt.Printf("processing %s\n", *src)
	if *src == "" || *out == "" {
		Usage()
		return
	}

	ProcessFile(*src, *out)
}

func FormatSource(filename string) error {
	// TODO: find gofmt dynamically.
	cmd := exec.Command("gofmt", "-w", filename)
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func AppendToFile(filename, text string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		return err
	}

	return nil
}
