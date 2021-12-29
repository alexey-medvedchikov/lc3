package bytecode

import (
	"fmt"
	"strconv"
)

type Register byte

const (
	R0 Register = 0
	R1 Register = 1
	R2 Register = 2
	R3 Register = 3
	R4 Register = 4
	R5 Register = 5
	R6 Register = 6
	R7 Register = 7
)

func (r *Register) String() string {
	return fmt.Sprintf("R%d", r)
}

func (r *Register) Capture(values []string) error {
	if len(values) != 1 {
		return fmt.Errorf("register can only capture single value: '%+v'", values)
	}

	v := values[0]
	i, err := strconv.ParseUint(v[1:], 10, 4)
	if err != nil {
		return err
	}
	*r = Register(i)

	return nil
}
