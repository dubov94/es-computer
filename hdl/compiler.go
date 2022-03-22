package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

const indent = "    "

type issue struct {
	exc   error
	token antlr.Token
}

type terminals struct {
	err   *issue
	name  string
	lower int
	upper int
}

func (self *terminals) String() string {
	return fmt.Sprintf("%s[%d:%d]", self.name, self.lower, self.upper)
}

type terminalsGroup struct {
	err   *issue
	items []*terminals
}

type connection struct {
	err    *issue
	source *terminals
	target *terminals
}

func (self *connection) String() string {
	return fmt.Sprintf("%s=%s", self.source.String(), self.target.String())
}

type connectionGroup struct {
	err   *issue
	items []*connection
}

type partImage struct {
	err         *issue
	name        string
	connections []*connection
}

func (part *partImage) String() string {
	var connections []string
	for _, connection := range part.connections {
		connections = append(connections, connection.String())
	}
	return fmt.Sprintf("%s(%s)", part.name, strings.Join(connections, ", "))
}

type partImageGroup struct {
	err   *issue
	items []*partImage
}

type chipImage struct {
	err     *issue
	name    string
	inputs  []*terminals
	outputs []*terminals
	parts   []*partImage
}

func (chip *chipImage) Lines() []string {
	var lines []string
	for _, input := range chip.inputs {
		lines = append(lines, fmt.Sprintf(`"%s" -> "%s";`, input.String(), chip.name))
	}
	for _, output := range chip.outputs {
		lines = append(lines, fmt.Sprintf(`"%s" -> "%s";`, chip.name, output.String()))
	}
	for _, part := range chip.parts {
		lines = append(lines, fmt.Sprintf(`"%s" -> "%s" [dir=none];`, chip.name, part.String()))
	}
	return lines
}

type hdlImage struct {
	err   *issue
	chips []*chipImage
}

func (hdl *hdlImage) Lines() []string {
	var lines []string
	lines = append(lines, "digraph hdlImage {")
	for _, chip := range hdl.chips {
		lines = append(lines, fmt.Sprintf("%ssubgraph {", indent))
		for _, line := range chip.Lines() {
			lines = append(lines, fmt.Sprintf("%s%s%s", indent, indent, line))
		}
		lines = append(lines, fmt.Sprintf("%s}", indent))
	}
	lines = append(lines, "}")
	return lines
}

type hdlVisitor struct {
	*BaseHdlVisitor
}

func (visitor *hdlVisitor) VisitSlice(context *SliceContext) interface{} {
	name := context.ID().GetText()
	fst := context.NUMBER(0)
	snd := context.NUMBER(1)
	if snd == nil && fst == nil {
		return &terminals{name: name, lower: 0, upper: 0}
	}
	if snd == nil {
		index, err := strconv.Atoi(fst.GetText())
		if err != nil {
			return &terminals{err: &issue{exc: err, token: fst.GetSymbol()}}
		}
		return &terminals{name: name, lower: index, upper: index + 1}
	}
	upper, err := strconv.Atoi(snd.GetText())
	if err != nil {
		return &terminals{err: &issue{exc: err, token: snd.GetSymbol()}}
	}
	lower, err := strconv.Atoi(fst.GetText())
	if err != nil {
		return &terminals{err: &issue{exc: err, token: fst.GetSymbol()}}
	}
	return &terminals{name: name, lower: lower, upper: upper + 1}
}

func (visitor *hdlVisitor) VisitConnection(context *ConnectionContext) interface{} {
	source := context.Slice(0).Accept(visitor).(*terminals)
	if source.err != nil {
		return &connection{err: source.err}
	}
	target := context.Slice(1).Accept(visitor).(*terminals)
	if target.err != nil {
		return &connection{err: target.err}
	}
	return &connection{source: source, target: target}
}

func (visitor *hdlVisitor) VisitConnections(context *ConnectionsContext) interface{} {
	var items []*connection
	for _, connectionContext := range context.AllConnection() {
		item := connectionContext.Accept(visitor).(*connection)
		if item.err != nil {
			return &connectionGroup{err: item.err}
		}
		items = append(items, item)
	}
	return &connectionGroup{items: items}
}

func (visitor *hdlVisitor) VisitPartDeclaration(context *PartDeclarationContext) interface{} {
	name := context.ID().GetText()
	connections := context.Connections().Accept(visitor).(*connectionGroup)
	if connections.err != nil {
		return &partImage{err: connections.err}
	}
	return &partImage{name: name, connections: connections.items}
}

func (visitor *hdlVisitor) VisitParts(context *PartsContext) interface{} {
	var items []*partImage
	for _, partContext := range context.AllPartDeclaration() {
		item := partContext.Accept(visitor).(*partImage)
		if item.err != nil {
			return &partImageGroup{err: item.err}
		}
		items = append(items, item)
	}
	return &partImageGroup{items: items}
}

func (visitor *hdlVisitor) VisitPinDeclaration(context *PinDeclarationContext) interface{} {
	name := context.ID().GetText()
	number := context.NUMBER()
	if number == nil {
		return &terminals{name: name, lower: 0, upper: 0}
	}
	upper, err := strconv.Atoi(number.GetText())
	if err != nil {
		return &terminals{err: &issue{exc: err, token: number.GetSymbol()}}
	}
	return &terminals{name: name, lower: 0, upper: upper}
}

func (visitor *hdlVisitor) VisitPinDeclarations(context *PinDeclarationsContext) interface{} {
	var items []*terminals
	for _, declarationContext := range context.AllPinDeclaration() {
		item := declarationContext.Accept(visitor).(*terminals)
		if item.err != nil {
			return &terminalsGroup{err: item.err}
		}
		items = append(items, item)
	}
	return &terminalsGroup{items: items}
}

func (visitor *hdlVisitor) VisitInputs(context *InputsContext) interface{} {
	return context.PinDeclarations().Accept(visitor)
}

func (visitor *hdlVisitor) VisitOutputs(context *OutputsContext) interface{} {
	return context.PinDeclarations().Accept(visitor)
}

func (visitor *hdlVisitor) VisitChip(context *ChipContext) interface{} {
	name := context.ID().GetText()
	inputGroup := context.Inputs().Accept(visitor).(*terminalsGroup)
	if inputGroup.err != nil {
		return &chipImage{err: inputGroup.err}
	}
	inputs := inputGroup.items
	outputGroup := context.Outputs().Accept(visitor).(*terminalsGroup)
	if outputGroup.err != nil {
		return &chipImage{err: outputGroup.err}
	}
	outputs := outputGroup.items
	partsGroup := context.Parts().Accept(visitor).(*partImageGroup)
	if partsGroup.err != nil {
		return &chipImage{err: partsGroup.err}
	}
	parts := partsGroup.items
	return &chipImage{name: name, inputs: inputs, outputs: outputs, parts: parts}
}

func (visitor *hdlVisitor) VisitChips(context *ChipsContext) interface{} {
	var chips []*chipImage
	for _, chipContext := range context.AllChip() {
		chips = append(chips, chipContext.Accept(visitor).(*chipImage))
	}
	return &hdlImage{chips: chips}
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
	fmt.Println(strings.Join(parser.Chips().Accept(visitor).(*hdlImage).Lines(), "\n"))
}
