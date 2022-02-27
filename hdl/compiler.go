package main

import (
    "fmt"
    "log"
    "os"

    "github.com/antlr/antlr4/runtime/Go/antlr"
)

type hdlImage struct {}

func (image *hdlImage) String() string {
    return "HdlImage"
}

type hdlVisitor struct {
    *BaseHdlVisitor
}

func (visitor *hdlVisitor) VisitChips(context *ChipsContext) interface{} {
    return &hdlImage{}
}

func main() {
    input, err := antlr.NewFileStream(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    lexer := NewHdlLexer(input)
    stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewHdlParser(stream)

    visitor := &hdlVisitor{}
    fmt.Println(parser.Chips().Accept(visitor))
}
