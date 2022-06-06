package main

import (
	"fmt"
	"os"

	"github.com/dubov94/es-computer/hdl/indexer"
	"github.com/dubov94/es-computer/hdl/reader"
)

func main() {
	hdlImage := reader.ReadHdl(os.Args[1])
	hdlIndex := indexer.IndexHdl(hdlImage)
	fmt.Println(hdlIndex)
}
