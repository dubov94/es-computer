package main

import (
    "fmt"
    "log"
    "os"

    "github.com/antlr/antlr4/runtime/Go/antlr"
)

type hdlListener struct {
    *BaseHdlListener
}

func (listener *hdlListener) EnterChip(context *ChipContext) {
    fmt.Println("EnterChip")
}

func (listener *hdlListener) ExitChip(context *ChipContext) {
    fmt.Println("ExitChip")
}

func main() {
    input, err := antlr.NewFileStream(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    lexer := NewHdlLexer(input)
    stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewHdlParser(stream)

    antlr.ParseTreeWalkerDefault.Walk(&hdlListener{}, parser.Chip())
}
