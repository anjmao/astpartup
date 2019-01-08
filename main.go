package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"

	"golang.org/x/tools/go/ast/astutil"
)

func valspec(name, typ string) *ast.ValueSpec {
	return &ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent(name)},
		Type: ast.NewIdent(typ),
	}
}

func callExpr() *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: "log",
			},
			Sel: &ast.Ident{
				Name: "Println",
			},
		},
		Args: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: `"new value"`,
			},
		},
	}
}

func main() {
	fset := token.NewFileSet()

	src := `
	package main

	import "fmt"

	func main() {
		fmt.Println("initial")
		var x int
	}
	`
	fast, err := parser.ParseFile(fset, "test.go", src, parser.AllErrors)
	if err != nil {
		log.Fatalf("could not parse file: %v", err)
	}

	fmt.Println("----ast before replace----")
	ast.Print(fset, fast)

	newCallExpr, err := parser.ParseExpr("log.Println(\"new value\")")
	if err != nil {
		log.Fatalf("could not parse expression: %v", err)
	}

	fmt.Println("----new expression----")
	ast.Print(token.NewFileSet(), newCallExpr)

	astutil.Apply(fast, func(c *astutil.Cursor) bool {
		switch c.Node().(type) {
		case *ast.ValueSpec:
			c.Replace(valspec("a", "int"))
			return false
		case *ast.CallExpr:
			// this works
			c.Replace(callExpr())

			// this doesn't work
			c.Replace(newCallExpr)
			return false
		}
		return true
	}, nil)

	fmt.Println("----ast after replace----")
	ast.Print(fset, fast)

	buf := &bytes.Buffer{}
	if err := format.Node(buf, fset, fast); err != nil {
		log.Fatalf("could not format new file: %v", err)
	}

	// print final result
	fmt.Println(buf.String())
}
