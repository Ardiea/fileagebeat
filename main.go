package main

import (
	"os"

	"github.com/Ardiea/fileagebeat/cmd"

	_ "github.com/Ardiea/fileagebeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
