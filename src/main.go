package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/pirespsps/synfo/parser"
)

func main() {
	if runtime.GOOS != "linux" {
		fmt.Print("Not in linux!")
		os.Exit(1)
	}

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
