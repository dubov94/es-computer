package reader

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type terminals struct {
	name  string
	lower int
	upper int
}

func (self *terminals) String() string {
	return fmt.Sprintf("%s[%d:%d]", self.name, self.lower, self.upper)
}

func (self *terminals) Name() string {
	return self.name
}

func (self *terminals) Lower() int {
	return self.lower
}

func (self *terminals) Upper() int {
	return self.upper
}

type connection struct {
	source *terminals
	target *terminals
}

func (self *connection) String() string {
	return fmt.Sprintf("%s=%s", self.target.String(), self.source.String())
}

func (self *connection) Source() *terminals {
	return self.source
}

func (self *connection) Target() *terminals {
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

type chipImage struct {
	name    string
	inputs  []*terminals
	outputs []*terminals
	parts   []*partImage
}

func (self *chipImage) String() string {
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

func (self *chipImage) Name() string {
	return self.name
}

func (self *chipImage) Inputs() []*terminals {
	return self.inputs
}

func (self *chipImage) Outputs() []*terminals {
	return self.outputs
}

func (self *chipImage) Parts() []*partImage {
	return self.parts
}

type hdlImage struct {
	chips []*chipImage
}

func (self *hdlImage) String() string {
	var chips []string
	for _, chip := range self.chips {
		chips = append(chips, chip.String())
	}
	return strings.Join(chips, "\n")
}

func (self *hdlImage) Chips() []*chipImage {
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
		return &terminals{name: name, lower: 0, upper: 0}
	}
	if snd == nil {
		index, err := strconv.Atoi(fst.GetText())
		if err != nil {
			visitor.reporter.HdlVisitorError(err, fst.GetSymbol())
		}
		return &terminals{name: name, lower: index, upper: index + 1}
	}
	upper, err := strconv.Atoi(snd.GetText())
	if err != nil {
		visitor.reporter.HdlVisitorError(err, snd.GetSymbol())
	}
	lower, err := strconv.Atoi(fst.GetText())
	if err != nil {
		visitor.reporter.HdlVisitorError(err, fst.GetSymbol())
	}
	return &terminals{name: name, lower: lower, upper: upper + 1}
}

func (visitor *hdlVisitor) VisitConnection(context *ConnectionContext) interface{} {
	target := context.Slice(0).Accept(visitor).(*terminals)
	source := context.Slice(1).Accept(visitor).(*terminals)
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
		return &terminals{name: name, lower: 0, upper: 0}
	}
	upper, err := strconv.Atoi(number.GetText())
	if err != nil {
		visitor.reporter.HdlVisitorError(err, number.GetSymbol())
	}
	return &terminals{name: name, lower: 0, upper: upper}
}

func (visitor *hdlVisitor) VisitPinDeclarations(context *PinDeclarationsContext) interface{} {
	var items []*terminals
	for _, declarationContext := range context.AllPinDeclaration() {
		item := declarationContext.Accept(visitor).(*terminals)
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
	inputs := context.Inputs().Accept(visitor).([]*terminals)
	outputs := context.Outputs().Accept(visitor).([]*terminals)
	parts := context.Parts().Accept(visitor).([]*partImage)
	return &chipImage{name: name, inputs: inputs, outputs: outputs, parts: parts}
}

func (visitor *hdlVisitor) VisitChips(context *ChipsContext) interface{} {
	var chips []*chipImage
	for _, chipContext := range context.AllChip() {
		chips = append(chips, chipContext.Accept(visitor).(*chipImage))
	}
	return &hdlImage{chips: chips}
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

func ReadHdl(path string) *hdlImage {
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
	return parser.Chips().Accept(visitor).(*hdlImage)
}
