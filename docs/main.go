package main

import (
	"log"

	"github.com/spf13/cobra/doc"
	"github.com/stevebargelt/dockhand/cmd"
)

func main() {

	err := doc.GenMarkdownTree(cmd.RootCmd, ".")
	if err != nil {
		log.Fatal(err)
	}
}
