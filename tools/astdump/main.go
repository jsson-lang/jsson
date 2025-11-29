package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"jsson/internal/ast"
	"jsson/internal/lexer"
	"jsson/internal/parser"
	"os"
)

func main() {
	in := flag.String("i", "", "input file")
	flag.Parse()
	if *in == "" {
		fmt.Println("provide -i file")
		os.Exit(1)
	}
	b, err := ioutil.ReadFile(*in)
	if err != nil {
		fmt.Println("read error:", err)
		os.Exit(1)
	}
	l := lexer.New(string(b))
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		fmt.Println("parser errors:")
		for _, e := range p.Errors() {
			fmt.Println("\t", e)
		}
		os.Exit(1)
	}

	for _, stmt := range prog.Statements {
		switch s := stmt.(type) {
		case *ast.AssignmentStatement:
			if s.Name.Value == "projects" {
				fmt.Println("Found assignment 'projects'")
				switch v := s.Value.(type) {
				case *ast.ArrayTemplate:
					fmt.Printf("Template keys=%v (count=%d)\n", v.Template.Keys, len(v.Template.Keys))
					for i, row := range v.Rows {
						fmt.Printf(" Row %d: cols=%d\n", i+1, len(row))
						for j, expr := range row {
							fmt.Printf("  - expr[%d] type=%T\n", j, expr)
						}
					}
					if v.Map != nil {
						fmt.Println("ArrayTemplate has Map clause; Map.Body keys:")
						for _, k := range v.Map.Body.Keys {
							expr := v.Map.Body.Properties[k]
							fmt.Printf(" - key=%s exprType=%T\n", k, expr)
						}
					}
				default:
					fmt.Printf("projects value is %T\n", v)
				}
			}
		}
	}
}
