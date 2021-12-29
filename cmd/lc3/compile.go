package main

import (
	"encoding/json"
	"github.com/alexey-medvedchikov/lc3/pkg/parser"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var compileCmd = func() cobra.Command {
	var outputFile string

	cmd := cobra.Command{
		Use:  "compile",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFile := args[0]
			return doCompile(inputFile, outputFile)
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "image.bin",
		"Output memory image")

	return cmd
}()

func doCompile(fPath string, outputFile string) error {
	inFp, err := os.OpenFile(fPath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := inFp.Close(); errClose != nil {
			log.Printf("[ERR] %s", errClose)
		}
	}()

	program, err := parser.Parse(inFp)
	if err != nil {
		return err
	}

	outFp, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := outFp.Close(); errClose != nil {
			log.Printf("[ERR] %s", errClose)
		}
	}()

	m := json.NewEncoder(outFp)
	m.SetIndent("", "  ")
	m.SetEscapeHTML(false)
	if err := m.Encode(program); err != nil {
		return err
	}

	return nil
}
