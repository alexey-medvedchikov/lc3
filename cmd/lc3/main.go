package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	rootCmd := &cobra.Command{
		Use: "lc3",
	}
	rootCmd.AddCommand(&runCmd)
	rootCmd.AddCommand(&compileCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
