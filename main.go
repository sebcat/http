package main

import (
	"fmt"
	"os"
)

type subCmdFunc func([]string) error

type subCmd struct {
	Name        string
	Description string
	Func        subCmdFunc
}

var subCmdList = []subCmd{
	subCmd{
		Name:        "mwu",
		Description: "Sample HTTP requests and perform the MW U test on two HTTP response-time groups",
		Func:        mwuMain,
	},

	subCmd{
		Name:        "file-server",
		Description: "Share a part of the local file system over HTTP",
		Func:        fsMain,
	},

	subCmd{
		Name:        "pong-server",
		Description: "Start an HTTP server that responds with \"pong\\n\"",
		Func:        psMain,
	},

	subCmd{
		Name:        "get-urls",
		Description: "Retrieve a list of HTTP resources and their status codes",
		Func:        guMain,
	},

	subCmd{
		Name:        "stress-test",
		Description: "Send HTTP requests at a specified rate and duration",
		Func:        stMain,
	},
}

func printSubCmds() {
	fmt.Println("Commands:")
	for _, cmd := range subCmdList {
		fmt.Printf("  http %v\n    %v\n", cmd.Name, cmd.Description)
	}
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" {
		printSubCmds()
		return
	}

	hasFunc := false
	for _, cmd := range subCmdList {
		if os.Args[1] == cmd.Name && cmd.Func != nil {
			hasFunc = true
			err := cmd.Func(os.Args[2:])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	if !hasFunc {
		fmt.Println("Command not implemented")
		os.Exit(1)
	}

	return
}
