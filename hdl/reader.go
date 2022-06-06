package reader

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type Terminals struct {
	name  string
	lower int
	upper int
}

func (self *Terminals) String() string {
	return fmt.Sprintf("%s[%d:%d]", self.name, self.lower, self.upper)
}

func (self *Terminals) Name() string {
	return self.name
}

func (self *Terminals) Lower() int {
	return self.lower
}

func (self *Terminals) Upper() int {
	return self.upper
}

type connection struct {
	source *Terminals
	target *Terminals
}

func (self *connection) String() string {
	return fmt.Sprintf("%s=%s", self.target.String(), self.source.String())
}

func (self *connection) Source() *Terminals {
	return self.source
}

func (self *connection) Target() *Terminals {
	return self.target
}

type partImage struct {
	name        string
	connections []*connection
}

func (self *partImage) String() string {
	var connections []string
	for _, connection := range self.connections {
		connections = append(connections, connection.String())
	}
	return fmt.Sprintf("%s(%s)", self.name, strings.Join(connections, ", "))
}

func (self *partImage) Name() string {
	return self.name
}

func (self *partImage) Connections() []*connection {
	return self.connections
}

type ChipImage struct {
	name    string
	inputs  []*Terminals
	outputs []*Terminals
	parts   []*partImage
}

func (self *ChipImage) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("CHIP %s", self.name))
	var inputs []string
	for _, input := range self.inputs {
		inputs = append(inputs, input.String())
	}
	lines = append(lines, fmt.Sprintf("IN %s;", strings.Join(inputs, ", ")))
	var outputs []string
	for _, output := range self.outputs {
		outputs = append(outputs, output.String())
	}
	lines = append(lines, fmt.Sprintf("OUT %s;", strings.Join(outputs, ", ")))
	lines = append(lines, "PARTS:")
	for _, part := range self.parts {
		lines = append(lines, fmt.Sprintf("%s;", part.String()))
	}
	return strings.Join(lines, "\n")
}

func (self *ChipImage) Name() string {
	return self.name
}

func (self *ChipImage) Inputs() []*Terminals {
	return self.inputs
}

func (self *ChipImage) Outputs() []*Terminals {
	return self.outputs
}

func (self *ChipImage) Parts() []*partImage {
	return self.parts
}

type HdlImage struct {
	chips []*ChipImage
}

func (self *HdlImage) String() string {
	var chips []string
	for _, chip := range self.chips {
		chips = append(chips, chip.String())
	}
	return strings.Join(chips, "\n")
}

func (self *HdlImage) Chips() []*ChipImage {
	return self.chips
}

type hdlVisitor struct {
	*BaseHdlVisitor
	reporter *errorListener
}

func (visitor *hdlVisitor) VisitSlice(context *SliceContext) interface{} {
	name := context.ID().GetText()
	fst := context.NUMBER(0)
	snd := context.NUMBER(1)
	if snd == nil && fst == nil {
		return &Terminals{name: name, lower: 0, upper: 0}
	}
	if snd == nil {
		index, err := strconv.Atoi(fst.GetText())
		if err != nil {
			visitor.reporter.HdlVisitorError(err, fst.GetSymbol())
		}
		return &Terminals{name: name, lower: index, upper: index + 1}
	}
	upper, err := strconv.Atoi(snd.GetText())
	if err != nil {
		visitor.reporter.HdlVisitorError(err, snd.GetSymbol())
	}
	lower, err := strconv.Atoi(fst.GetText())
	if err != nil {
		visitor.reporter.HdlVisitorError(err, fst.GetSymbol())
	}
	return &Terminals{name: name, lower: lower, upper: upper + 1}
}

func (visitor *hdlVisitor) VisitConnection(context *ConnectionContext) interface{} {
	target := context.Slice(0).Accept(visitor).(*Terminals)
	source := context.Slice(1).Accept(visitor).(*Terminals)
	return &connection{source: source, target: target}
}

func (visitor *hdlVisitor) VisitConnections(context *ConnectionsContext) interface{} {
	var items []*connection
	for _, connectionContext := range context.AllConnection() {
		item := connectionContext.Accept(visitor).(*connection)
		items = append(items, item)
	}
	return items
}

func (visitor *hdlVisitor) VisitPartDeclaration(context *PartDeclarationContext) interface{} {
	name := context.ID().GetText()
	connections := context.Connections().Accept(visitor).([]*connection)
	return &partImage{name: name, connections: connections}
}

func (visitor *hdlVisitor) VisitParts(context *PartsContext) interface{} {
	var items []*partImage
	for _, partContext := range context.AllPartDeclaration() {
		item := partContext.Accept(visitor).(*partImage)
		items = append(items, item)
	}
	return items
}

func (visitor *hdlVisitor) VisitPinDeclaration(context *PinDeclarationContext) interface{} {
	name := context.ID().GetText()
	number := context.NUMBER()
	if number == nil {
		return &Terminals{name: name, lower: 0, upper: 0}
	}
	upper, err := strconv.Atoi(number.GetText())
	if err != nil {
		visitor.reporter.HdlVisitorError(err, number.GetSymbol())
	}
	return &Terminals{name: name, lower: 0, upper: upper}
}

func (visitor *hdlVisitor) VisitPinDeclarations(context *PinDeclarationsContext) interface{} {
	var items []*Terminals
	for _, declarationContext := range context.AllPinDeclaration() {
		item := declarationContext.Accept(visitor).(*Terminals)
		items = append(items, item)
	}
	return items
}

func (visitor *hdlVisitor) VisitInputs(context *InputsContext) interface{} {
	return context.PinDeclarations().Accept(visitor)
}

func (visitor *hdlVisitor) VisitOutputs(context *OutputsContext) interface{} {
	return context.PinDeclarations().Accept(visitor)
}

func (visitor *hdlVisitor) VisitChip(context *ChipContext) interface{} {
	name := context.ID().GetText()
	inputs := context.Inputs().Accept(visitor).([]*Terminals)
	outputs := context.Outputs().Accept(visitor).([]*Terminals)
	parts := context.Parts().Accept(visitor).([]*partImage)
	return &ChipImage{name: name, inputs: inputs, outputs: outputs, parts: parts}
}

func (visitor *hdlVisitor) VisitChips(context *ChipsContext) interface{} {
	var chips []*ChipImage
	for _, chipContext := range context.AllChip() {
		chips = append(chips, chipContext.Accept(visitor).(*ChipImage))
	}
	return &HdlImage{chips: chips}
}

type errorListener struct {
	*antlr.DefaultErrorListener
}

func (listener *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	log.Fatalf("[%d:%d] %s", line, column, msg)
}

func (listener *errorListener) HdlVisitorError(exc error, token antlr.Token) {
	log.Fatalf("[%d:%d] %s", token.GetLine(), token.GetColumn(), exc.Error())
}

func ReadHdl(path string) *HdlImage {
	input, err := antlr.NewFileStream(path)
	if err != nil {
		log.Fatal(err)
	}

	reporter := &errorListener{}

	lexer := NewHdlLexer(input)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(reporter)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewHdlParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(reporter)

	visitor := &hdlVisitor{reporter: reporter}
	return parser.Chips().Accept(visitor).(*HdlImage)
}
