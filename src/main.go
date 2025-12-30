package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"slices"

	"github.com/pirespsps/synfo/parser"
)

var cmds = []string{
	"storage",
	"cpu",
	"ram",
	"graphics",
	"hardware",
	"network",
	"system",
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Print("Not in linux!")
		os.Exit(1)
	}

	var isJson bool

	flag.BoolVar(&isJson, "J", false, "Set the output to JSON format")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Print("missing command\n")
		os.Exit(1)
	}

	cmd := args[0]
	if !slices.Contains(cmds, cmd) {
		fmt.Print("option doesn't exist\n")
		os.Exit(1)
	}

	option := "overall" //overall,extensive,monitor
	if len(args) > 1 {
		option = args[1]
	}

	resp, err := parser.GetResponse(cmd, option)
	if err != nil {
		panic(err)
	}

	if isJson {
		if err := resp.PrintJson(); err != nil {
			panic(err)
		}
	} else {
		resp.Print()
	}
}
