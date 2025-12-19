package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pirespsps/synfo/parser"
)

func main() {
	// args, with cpu, memory, network, all, etc.....
	//arg with json input option
	var isJson bool

	flag.BoolVar(&isJson, "j", false, "Set the output to JSON format")

	flag.Parse()

	option := os.Args[1]

	data, err := parser.FetchData(option, isJson)
	if err != nil {
		panic(err)
	}

	fmt.Print(data)

}
