package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"jsson/internal/lexer"
	"jsson/internal/token"
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
	l.SetSourceFile(*in)
	for {
		t := l.NextToken()
		fmt.Printf("%4d:%-3d %-15s %q\n", t.Line, t.Column, string(t.Type), t.Literal)
		if t.Type == token.EOF {
			break
		}
	}
}
