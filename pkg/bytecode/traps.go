package bytecode

const (
	KeyboardStatusReg = 0xFE00
	KeyboardDataReg   = 0xFE02

	DisplayStatusReg = 0xFE04
	DisplayDataReg   = 0xFE06

	ControlReg = 0xFFFE
)

const (
	TrapGETCAddr  uint16 = 0x0020
	TrapOUTAddr   uint16 = 0x0021
	TrapPUTSAddr  uint16 = 0x0022
	TrapINAddr    uint16 = 0x0023
	TrapPUTSPAddr uint16 = 0x0024
	TrapHALTAddr  uint16 = 0x0025
)

var (
	TrapGETC = []uint16{
		Trap(0x25),
		RET(),
	}

	TrapOUT = []uint16{
		Trap(0x25),
		RET(),
	}

	TrapPUTS = []uint16{
		STR(R1, R6, 0), // push R1
		AddImm(R6, R6, -1),
		STR(R2, R6, 0), // push R2
		AddImm(R6, R6, -1),

		// NUL-string is in R0
		// LOOP:
		LDI(R1, 15),        // R1 <- mem[DisplayStatusReg]
		LD(R2, 16),         // R2 <- 0b1000_0000_0000_0000
		AndReg(R1, R1, R2), // R1 <- R1 & R2
		BRx(0b010, -4),
		LDR(R1, R0, 0),     // R1 <- mem[R0]
		BRx(0b010, 9),      // goto RETURN if loaded 0x0000
		LD(R2, 12),         // R2 <- 0x00FF
		AndReg(R1, R1, R2), // R1 <- R1 & R2
		STI(R1, 8),         // mem[DisplayDataReg] <- R1
		AddImm(R0, R0, 1),  // R0 <- R0 + 1
		BRx(0b111, -10),    // goto LOOP

		AddImm(R6, R6, 1), // pop R2
		LDR(R2, R6, 0),
		AddImm(R6, R6, 1), // pop R1
		LDR(R1, R6, 0),

		// RETURN:
		RET(),

		DisplayStatusReg,
		DisplayDataReg,
		0b1000_0000_0000_0000,
		0x00FF,
	}

	TrapIN = []uint16{
		Trap(0x25),
		RET(),
	}

	TrapPUTSP = []uint16{
		Trap(0x25),
		RET(),
	}

	TrapHALT = []uint16{
		STR(R0, R6, 0), // push R0
		AddImm(R6, R6, -1),
		STR(R1, R6, 0), // push R1
		AddImm(R6, R6, -1),

		LD(R0, 9),          // R0 <- 0b0111_1111_1111_1111
		LDI(R1, 7),         // R1 <- mem[ControlReg]
		AndReg(R1, R1, R0), // R1 <- R1 & R0
		STI(R1, 5),         // mem[ControlReg] <- R1

		AddImm(R6, R6, 1), // pop R1
		LDR(R1, R6, 0),
		AddImm(R6, R6, 1), // pop R0
		LDR(R0, R6, 0),
		RET(),
		ControlReg,
		0b0111_1111_1111_1111,
	}
)
