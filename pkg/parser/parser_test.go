package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const code = `        .ORIG x0000

; the TRAP vector table
    .FILL BAD_TRAP  ; x00
    .FILL BAD_TRAP  ; x01
; the interrupt vector table
; interrupts are not currently implemented


;;; OS_START - operating system entry point (always starts at x0200)
OS_START EXTRA_LABEL
    ;; set MPR
    LD R0, MPR_INIT
    STI R0, OS_MPR

    ;; set timer interval
    LD R0, TIM_INIT
    STI R0, OS_TMI

    ;; start running user code (clear Privilege bit w/ JMPT)
    LD R7, USER_CODE_ADDR
    JMPT    R7

OS_KBSR .FILL xFE00     ; keyboard status register
OS_KBDR .FILL xFE02     ; keyboard data register

OS_SAVE_R0      .BLKW 1
OS_SAVE_R1      .BLKW 1

MASK_HI         .FILL x7FFF
TIM_INIT        .FILL #40
MPR_INIT    .FILL xFFFF ; user can access everything
;MPR_INIT   .FILL x0FF8 ; user can access x3000 to xbfff

;;; GETC - Read a single character of input from keyboard device into R0
TRAP_GETC
    LDI R0,OS_KBSR      ; wait for a keystroke
    BRzp TRAP_GETC
    LDI R0,OS_KBDR      ; read it and return
    RET


;;; OUT - Write the character in R0 to the console.
TRAP_OUT
    ST R1,OS_OUT_SAVE_R1    ; save R1
TRAP_OUT_WAIT
    LDI R1,OS_DSR       ; wait for the display to be ready
    BRzp TRAP_OUT_WAIT
    STI R0,OS_DDR       ; write the character and return
    LD R1,OS_OUT_SAVE_R1    ; restore R1
    RET


;;; PUTS - Write a NUL-terminated string of characters to the console,
;;; starting from the address in R0.
TRAP_PUTS
    ST R0,OS_SAVE_R0    ; save R0, R1, and R7
    ST R1,OS_SAVE_R1
    ST R7,OS_SAVE_R7
    ADD R1,R0,#0        ; move string pointer (R0) into R1

TRAP_PUTS_LOOP
    LDR R0,R1,#0        ; write characters in string using OUT
    BRz TRAP_PUTS_DONE
    OUT
    ADD R1,R1,#1
    BRnzp TRAP_PUTS_LOOP

TRAP_PUTS_DONE
    LD R0,OS_SAVE_R0    ; restore R0, R1, and R7
    LD R1,OS_SAVE_R1
    LD R7,OS_SAVE_R7
    RET

;;; IN - prompt the user for a single character input, which is stored
;;; in R0 and also echoed to the console.
TRAP_IN
    ST R7,OS_IN_SAVE_R7 ; save R7 (no need to save R0, since we
                ;    overwrite later
    LEA R0,TRAP_IN_MSG  ; prompt for input
    PUTS
    GETC            ; read a character
    OUT         ; echo back to monitor
    ST R0,OS_SAVE_R0    ; save the character
    AND R0,R0,#0        ; write a linefeed, too
    ADD R0,R0,#10
    OUT
    LD R0,OS_SAVE_R0    ; restore the character
    LD R7,OS_IN_SAVE_R7 ; restore R7
    RET                     ; this doesn't work, because


;;; PUTSP - Write a NUL-terminated string of characters, packed 2 per
;;; memory location, to the console, starting from the address in R0.
TRAP_PUTSP
    ; NOTE: This trap will end when it sees any NUL, even in
    ; packed form, despite the P&P second edition's requirement
    ; of a double NUL.

    ST R0,OS_SAVE_R0    ; save R0, R1, R2, R3, and R7
    ST R1,OS_SAVE_R1
    ST R2,OS_SAVE_R2
    ST R3,OS_SAVE_R3
    ST R7,OS_SAVE_R7
    ADD R1,R0,#0        ; move string pointer (R0) into R1

TRAP_PUTSP_LOOP
    LDR R2,R1,#0        ; read the next two characters
    LD R0,LOW_8_BITS    ; use mask to get low byte
    AND R0,R0,R2        ; if low byte is NUL, quit printing
    BRz TRAP_PUTSP_DONE
    OUT         ; otherwise print the low byte

    AND R0,R0,#0        ; shift high byte into R0
    ADD R3,R0,#8
TRAP_PUTSP_S_LOOP
    ADD R0,R0,R0        ; shift R0 left
    ADD R2,R2,#0        ; move MSB from R2 into R0
    BRzp TRAP_PUTSP_MSB_0
    ADD R0,R0,#1
TRAP_PUTSP_MSB_0
    ADD R2,R2,R2        ; shift R2 left
    ADD R3,R3,#-1
    BRp TRAP_PUTSP_S_LOOP

    ADD R0,R0,#0        ; if high byte is NUL, quit printing
    BRz TRAP_PUTSP_DONE
    OUT         ; otherwise print the low byte

    ADD R1,R1,#1        ; and keep going
    BRnzp TRAP_PUTSP_LOOP

TRAP_PUTSP_DONE
    LD R0,OS_SAVE_R0    ; restore R0, R1, R2, R3, and R7
    LD R1,OS_SAVE_R1
    LD R2,OS_SAVE_R2
    LD R3,OS_SAVE_R3
    LD R7,OS_SAVE_R7
    RET


;;; HALT - trap handler for halting machine
TRAP_HALT
    LDI R0,OS_MCR
    LD R1,MASK_HI           ; clear run bit in MCR
    AND R0,R0,R1
    STI R0,OS_MCR       ; halt!
    BRnzp OS_START      ; restart machine


;;; BAD_TRAP - code to execute for undefined trap
BAD_TRAP
    BRnzp TRAP_HALT     ; execute HALT


;;; BAD_INT - code to execute for undefined interrupt. There won't
;;; actually be any interrupts, so this will never actually get called.
BAD_INT     RTI

TRAP_IN_MSG .STRINGZ "\nInput a character> "`

func TestParse(t *testing.T) {

	in := bytes.NewBuffer([]byte(code))

	_, err := Parse(in)
	assert.NoError(t, err)
}
