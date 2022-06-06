package indexer

import (
	"fmt"
	"log"
	"strings"

	"github.com/dubov94/es-computer/hdl/reader"
)

type PortKind int

const (
	portKindInput PortKind = iota
	portKindOutput
)

func (kind PortKind) IsInput() bool {
	return kind == portKindInput
}

func (kind PortKind) IsOutput() bool {
	return kind == portKindOutput
}

type Port struct {
	kind  PortKind
	lower int
	upper int
}

func (self *Port) String() string {
	return fmt.Sprintf("(%d, [%d:%d])", self.kind, self.lower, self.upper)
}

type PortMapping map[string]*Port

func (self *PortMapping) String() string {
	var pairs []string
	for name, port := range *self {
		pairs = append(pairs, fmt.Sprintf("%s: %s", name, port.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

type ChipIndex struct {
	nameToPort PortMapping
}

func (self *ChipIndex) String() string {
	return fmt.Sprintf("{PortMapping: %s}", self.nameToPort.String())
}

type ChipMapping map[string]*ChipIndex

func (self *ChipMapping) String() string {
	var pairs []string
	for name, chipIndex := range *self {
		pairs = append(pairs, fmt.Sprintf("%s: %s", name, chipIndex.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

type HdlIndex struct {
	nameToChip ChipMapping
}

func (self *HdlIndex) String() string {
	return fmt.Sprintf("{ChipMapping: %s}", self.nameToChip.String())
}

func (self PortMapping) ingest(kind PortKind, terminals *reader.Terminals) error {
	name := terminals.Name()
	if _, isPresent := self[name]; isPresent {
		return fmt.Errorf("'%s' redeclaration", name)
	}
	self[name] = &Port{
		kind:  kind,
		lower: terminals.Lower(),
		upper: terminals.Upper(),
	}
	return nil
}

func indexChip(chipImage *reader.ChipImage) (*ChipIndex, error) {
	var err error
	chipName := chipImage.Name()
	nameToPort := make(PortMapping)
	for _, terminals := range chipImage.Inputs() {
		if err = nameToPort.ingest(portKindInput, terminals); err != nil {
			return nil, fmt.Errorf("%s: %v", chipName, err)
		}
	}
	for _, terminals := range chipImage.Outputs() {
		if err = nameToPort.ingest(portKindOutput, terminals); err != nil {
			return nil, fmt.Errorf("%s: %v", chipName, err)
		}
	}
	return &ChipIndex{nameToPort}, nil
}

func (self ChipMapping) ingest(chipImage *reader.ChipImage) error {
	chipName := chipImage.Name()
	if _, isPresent := self[chipName]; isPresent {
		return fmt.Errorf("%s: `CHIP` redeclartion", chipName)
	}
	chipIndex, err := indexChip(chipImage)
	if err != nil {
		return err
	}
	self[chipName] = chipIndex
	return nil
}

func IndexHdl(hdlImage *reader.HdlImage) *HdlIndex {
	var err error
	nameToChip := make(ChipMapping)
	for _, chipImage := range hdlImage.Chips() {
		if err = nameToChip.ingest(chipImage); err != nil {
			log.Fatalf("%v", err)
		}
	}
	return &HdlIndex{nameToChip}
}
