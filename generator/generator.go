package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"regexp"
	"strings"

	"github.com/matematik7/counterfeiter/model"

	"golang.org/x/tools/imports"
)

type CodeGenerator struct {
	Model       model.InterfaceToFake
	StructName  string
	PackageName string
}

func (gen CodeGenerator) GenerateFake() (string, error) {
	buf := new(bytes.Buffer)
	err := format.Node(buf, token.NewFileSet(), gen.sourceFile())
	if err != nil {
		return "", err
	}

	code, err := imports.Process("", buf.Bytes(), nil)
	return commentLine() + prettifyCode(string(code)), err
}

func (gen CodeGenerator) sourceFile() ast.Node {
	declarations := []ast.Decl{
		gen.imports(),
		gen.fakeStructType(),
	}

	for _, method := range gen.Model.Methods {
		declarations = append(
			declarations,
			gen.methodImplementation(method),
		)
	}

	return &ast.File{
		Name:  &ast.Ident{Name: gen.PackageName},
		Decls: declarations,
	}
}

func (gen CodeGenerator) imports() ast.Decl {
	specs := []ast.Spec{
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"` + gen.Model.ImportPath + `"`,
			},
		},
	}

	for _, spec := range gen.Model.ImportSpecs {
		specs = append(specs, spec)
	}

	return &ast.GenDecl{
		Lparen: 1,
		Tok:    token.IMPORT,
		Specs:  specs,
	}
}

func (gen CodeGenerator) fakeStructType() ast.Decl {
	structFields := []*ast.Field{}

	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: gen.StructName},
				Type: &ast.StructType{
					Fields: &ast.FieldList{List: structFields},
				},
			},
		},
	}
}

func (gen CodeGenerator) methodImplementation(method *ast.Field) *ast.FuncDecl {
	methodType := method.Type.(*ast.FuncType)

	var lastStatement ast.Stmt
	if methodType.Results.NumFields() > 0 {
		returnValues := []ast.Expr{}
		eachMethodResult(methodType, func(name string, t ast.Expr) {
			returnValues = append(returnValues, &ast.BasicLit{
				Kind:  token.STRING,
				Value: "nil",
			})
		})

		lastStatement = &ast.ReturnStmt{Results: returnValues}
	} else {
		lastStatement = &ast.ReturnStmt{Results: nil}
	}

	return &ast.FuncDecl{
		Name: method.Names[0],
		Type: &ast.FuncType{
			Params:  methodType.Params,
			Results: methodType.Results,
		},
		Recv: gen.receiverFieldList(),
		Body: &ast.BlockStmt{List: []ast.Stmt{
			lastStatement,
		}},
	}
}

func (gen CodeGenerator) receiverFieldList() *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{receiverIdent()},
				Type:  &ast.StarExpr{X: ast.NewIdent(gen.StructName)},
			},
		},
	}
}

func eachMethodParam(methodType *ast.FuncType, cb func(string, ast.Expr, int)) {
	i := 0
	for _, field := range methodType.Params.List {
		if len(field.Names) == 0 {
			cb(fmt.Sprintf("arg%d", i+1), field.Type, i)
			i++
		} else {
			for _, name := range field.Names {
				cb(name.Name, field.Type, i)
				i++
			}
		}
	}
}

func eachMethodResult(methodType *ast.FuncType, cb func(string, ast.Expr)) {
	for i, field := range methodType.Results.List {
		cb(fmt.Sprintf("result%d", i+1), field.Type)
	}
}

func argsStructTypeForMethod(methodType *ast.FuncType) *ast.StructType {
	fields := []*ast.Field{}

	eachMethodParam(methodType, func(name string, t ast.Expr, _ int) {
		fields = append(fields, &ast.Field{
			Type:  storedTypeForType(t),
			Names: []*ast.Ident{ast.NewIdent(name)},
		})
	})

	return &ast.StructType{
		Fields: &ast.FieldList{List: fields},
	}
}

func returnStructTypeForMethod(methodType *ast.FuncType) *ast.StructType {
	resultFields := []*ast.Field{}
	eachMethodResult(methodType, func(name string, t ast.Expr) {
		resultFields = append(resultFields, &ast.Field{
			Type:  t,
			Names: []*ast.Ident{ast.NewIdent(name)},
		})
	})

	return &ast.StructType{
		Fields: &ast.FieldList{List: resultFields},
	}
}

func storedTypeForType(t ast.Expr) ast.Expr {
	if ellipsis, ok := t.(*ast.Ellipsis); ok {
		return &ast.ArrayType{Elt: ellipsis.Elt}
	} else {
		return t
	}
}

func callCountMethodName(method *ast.Field) string {
	return method.Names[0].Name + "CallCount"
}

func callArgsMethodName(method *ast.Field) string {
	return method.Names[0].Name + "ArgsForCall"
}

func callArgsFieldName(method *ast.Field) string {
	return privatize(callArgsMethodName(method))
}

func mutexFieldName(method *ast.Field) string {
	return privatize(method.Names[0].Name) + "Mutex"
}

func methodStubFuncName(method *ast.Field) string {
	return method.Names[0].Name + "Stub"
}

func returnSetterMethodName(method *ast.Field) string {
	return method.Names[0].Name + "Returns"
}

func returnStructFieldName(method *ast.Field) string {
	return privatize(returnSetterMethodName(method))
}

func receiverIdent() *ast.Ident {
	return ast.NewIdent("fake")
}

func callMutex(method *ast.Field, verb string) ast.Stmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X:   receiverIdent(),
					Sel: ast.NewIdent(mutexFieldName(method)),
				},
				Sel: ast.NewIdent(verb),
			},
		},
	}
}

func deferMutex(method *ast.Field, verb string) ast.Stmt {
	return &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X:   receiverIdent(),
					Sel: ast.NewIdent(mutexFieldName(method)),
				},
				Sel: ast.NewIdent(verb),
			},
		},
	}
}

func publicize(input string) string {
	return strings.ToUpper(input[0:1]) + input[1:]
}

func privatize(input string) string {
	return strings.ToLower(input[0:1]) + input[1:]
}

func nilCheck(x ast.Expr) ast.Expr {
	return &ast.BinaryExpr{
		X:  x,
		Op: token.NEQ,
		Y: &ast.BasicLit{
			Kind:  token.STRING,
			Value: "nil",
		},
	}
}

func commentLine() string {
	return "// This file was generated by counterfeiter\n"
}

func prettifyCode(code string) string {
	code = funcRegexp.ReplaceAllString(code, "\n\nfunc")
	code = emptyStructRegexp.ReplaceAllString(code, "struct{}")
	code = strings.Replace(code, "\n\n\n", "\n\n", -1)
	return code
}

var funcRegexp = regexp.MustCompile("\nfunc")
var emptyStructRegexp = regexp.MustCompile("struct[\\s]+{[\\s]+}")
