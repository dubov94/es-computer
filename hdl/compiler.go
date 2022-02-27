package main

import (
    "fmt"
    "log"
    "os"

    "github.com/antlr/antlr4/runtime/Go/antlr"
)

type hdlVisitor struct {
    *BaseHdlVisitor
}

type hdlImage struct {}

func (visitor *hdlVisitor) VisitChips(context *ChipsContext) interface{} {
    fmt.Println(context.GetText())
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
    parser.Chips().Accept(visitor)
}
