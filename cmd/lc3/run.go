package main

import (
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/alexey-medvedchikov/lc3/pkg/machine"
)

var runCmd = func() cobra.Command {
	var startAddr uint16
	var enableTrace bool

	cmd := cobra.Command{
		Use:  "run",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imagePath := args[0]

			return doRun(imagePath, startAddr, enableTrace)
		},
	}

	cmd.Flags().Uint16VarP(&startAddr, "start-addr", "s", machine.UserStart,
		"Initial Program Counter value")
	cmd.Flags().BoolVarP(&enableTrace, "trace", "t", false, "Enable tracing")

	return cmd
}()

func doRun(imagePath string, startAddr uint16, enableTrace bool) error {
	var m machine.Machine

	fp, err := os.OpenFile(imagePath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := fp.Close(); errClose != nil {
			log.Printf("[ERR] %s", errClose)
		}
	}()

	w := machine.NewMemoryWriter(&m.Memory)
	written, err := io.Copy(w, fp)
	log.Printf("[INFO] %d bytes written", written)
	if err != nil {
		if !errors.Is(err, io.ErrShortWrite) {
			return err
		}
		log.Printf("[WARN] input is too big, truncated to memory size")
	}

	log.Printf("[INFO] starting VM (PC = 0x%0.4x)", startAddr)

	m.Init()
	var traceFunc func(m *machine.Machine)
	if enableTrace {
		traceFunc = func(m *machine.Machine) {
			time.Sleep(1 * time.Second)
		}
		t := machine.NewTracedMachine(os.Stdout, &m)
		t.Start(traceFunc)
	} else {
		m.Start()
	}

	return nil
}
