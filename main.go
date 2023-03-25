package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

var tabCounter int = 0

func convertGoType(t string) string {
	switch t {
	case "int":
		return t
	case "string":
		return "char*"
	case "bool":
		return "int"
	case "uint8", "uint16", "uint32", "uint64":
		return "uint" + strings.Replace(t, "uint", "", 1) + "_t"
	default:
		fmt.Println("Warning: could not convert type", t)
		return t
	}
}

func compile(inputCode string) string {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", inputCode, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	compiledCode := ""

	ast.Inspect(f, func(node ast.Node) bool {
		compiledCode += switchNode(node)
		return true
	})

	return compiledCode
}

func tokenHandle(node token.Token) string {
	return node.String()
}

func exprHandle(node ast.Expr) string {
	switch node := node.(type) {
	case *ast.Ident:
		return node.Name
	case *ast.BasicLit:
		return node.Value
	case *ast.BinaryExpr:
		return fmt.Sprintf("%s %s %s", exprHandle(node.X), tokenHandle(node.Op), exprHandle(node.Y))
	case *ast.CallExpr:
		args := make([]string, 0)
		for _, arg := range node.Args {
			args = append(args, exprHandle(arg))
		}
		return fmt.Sprintf("%s(%s)", node.Fun.(*ast.Ident).Name, strings.Join(args, ", "))
	default:
		return "ye"
	}
}

func stmtHandle(node ast.Stmt) string {
	switch node := node.(type) {
	case *ast.DeclStmt:
		xVar := node.Decl.(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
		varName := xVar.Names[0]
		varType := convertGoType(fmt.Sprint(xVar.Type))
		varVal := exprHandle(xVar.Values[0])

		return fmt.Sprintf("%s %s = %s;\n", varType, varName, varVal)
	case *ast.ReturnStmt:
		return fmt.Sprintf("return %s;\n", exprHandle(node.Results[0]))
	case *ast.ExprStmt:
		return exprHandle(node.X) + ";\n"
	case *ast.IfStmt:
		output := fmt.Sprintf("if (%s)\n%s{\t\n", exprHandle(node.Cond), strings.Repeat("\t", tabCounter))

		tabCounter++
		for _, stmt := range node.Body.List {
			output += strings.Repeat("\t", tabCounter) + stmtHandle(stmt)
		}
		tabCounter--

		output += strings.Repeat("\t", tabCounter) + "}\n"

		// Check for else
		if len(node.Else.(*ast.BlockStmt).List) != 0 {
			output += strings.Repeat("\t", tabCounter) + "else\n" + strings.Repeat("\t", tabCounter) + "{\n"
			tabCounter++
			for _, stmt := range node.Else.(*ast.BlockStmt).List {
				output += strings.Repeat("\t", tabCounter) + stmtHandle(stmt)
			}
			tabCounter--
			output += strings.Repeat("\t", tabCounter) + "}\n"
		}

		return output

	default:
		return "\n"
	}
}

func switchNode(node ast.Node) string {
	compiledBit := ""
	switch node := node.(type) {
	case *ast.ImportSpec:
		compiledBit += fmt.Sprintf("#include %s\n", node.Path.Value)

	case *ast.FuncDecl:
		functionType := convertGoType(fmt.Sprint(node.Type.Results.List[0].Type))
		functionParameters := make([]string, 0)

		for _, param := range node.Type.Params.List {
			for i := range param.Names {
				functionParameters = append(functionParameters, fmt.Sprintf("%s %s", convertGoType(node.Type.Params.List[i].Type.(*ast.Ident).Name), param.Names[i].Name))
			}
		}

		compiledBit += fmt.Sprintf("\n%s %s(%s)\n{\n", functionType, node.Name, strings.Join(functionParameters, ", "))

		tabCounter++
		for _, stmt := range node.Body.List {
			compiledBit += strings.Repeat("\t", tabCounter) + stmtHandle(stmt)
		}
		tabCounter--

		compiledBit += "}\n"
	default:
		// Unhandled token types
	}

	return compiledBit
}

func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Printf("USAGE: %s <input_file> <output_file>", args[0])
		return
	}

	inputFile := args[1]
	outputFile := args[2]

	sourceCode, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error: Could not open file '%s'\n", inputFile)
		fmt.Printf("Reason: %s", err)
		return
	}

	compiledCode := compile(string(sourceCode))

	err = os.WriteFile(outputFile, []byte(compiledCode), 0644)

	if err != nil {
		fmt.Printf("Error: Could not write compiled code to file '%s'\n", outputFile)
		fmt.Printf("Reason: %s", err)
		return
	}
}
