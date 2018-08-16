package main

import (
	"fmt"
)

type CPU struct {
	PC     uint16
	SP     uint8
	A      uint8
	X      uint8
	Y      uint8
	F      byte // 8 bit status flags
	CFlag  bool
	ZFlag  bool
	IFlag  bool
	DFlag  bool
	BFlag  bool
	UFlag  bool // unused
	VFlag  bool
	NFlag  bool
	Cycles uint
	Memory [100000]byte
	Debug  bool
}

func (cpu *CPU) Print() {
	fmt.Printf("PC: %d SP: %d A: %d X: %d Y: %d, Cycles: %d\n", cpu.PC, cpu.SP, cpu.A, cpu.X, cpu.Y, cpu.Cycles)
}

func (cpu *CPU) Exec() {
	opcode := cpu.Memory[cpu.PC]
	context := context(cpu, opcode)
	instruction := instructionMap[opcode]
	metadata := instructionMetadata[opcode]

	if cpu.Debug {
		fmt.Println("Beginning execution...")
		fmt.Printf("Opcode: %#x, Addressing Mode: %d, Address: %d\n", opcode, context.Mode, context.Address)
		cpu.Print()
	}

	cpu.PC += metadata.Bytes
	instruction.Exec(cpu, context)
	cpu.Cycles += metadata.Cycles(context)

	if cpu.Debug {
		fmt.Println("Finished execution")
		cpu.Print()
		fmt.Println()
	}
}

func (cpu *CPU) setZeroFlag(n byte) {
	// set zero flag if input is zero
	cpu.ZFlag = n == 0
}

func (cpu *CPU) setNegativeFlag(n byte) {
	// Set negative flag if input bit 7 is 1
	cpu.NFlag = n&0x80 > 0
}

func (cpu *CPU) setZeroAndNegativeFlags(n byte) {
	cpu.setZeroFlag(n)
	cpu.setNegativeFlag(n)
}
func (cpu *CPU) flagToInt(flag bool) uint8 {
	if flag {
		return 1
	} else {
		return 0
	}
}

func (cpu *CPU) intToFlag(n uint8) bool {
	if n == 0 {
		return false
	} else {
		return true
	}
}

func (cpu *CPU) flagsToByte() byte {
	return (cpu.flagToInt(cpu.NFlag) << 7) |
		(cpu.flagToInt(cpu.VFlag) << 6) |
		(cpu.flagToInt(cpu.UFlag) << 5) |
		(cpu.flagToInt(cpu.BFlag) << 4) |
		(cpu.flagToInt(cpu.DFlag) << 3) |
		(cpu.flagToInt(cpu.IFlag) << 2) |
		(cpu.flagToInt(cpu.ZFlag) << 1) |
		(cpu.flagToInt(cpu.CFlag) << 0)
}

func (cpu *CPU) stackPush(value byte) {
	cpu.Memory[0x100|uint16(cpu.SP)] = value
	cpu.SP -= 1
}

func (cpu *CPU) stackPop() byte {
	cpu.SP += 1
	return cpu.Memory[0x100|uint16(cpu.SP)]
}

func context(cpu *CPU, opcode byte) *InstructionContext {
	var address uint16

	pageCrossed := false
	mode := addressingModeMap[opcode]

	if mode == Immediate {
		address = cpu.PC + 1
	} else if mode == ZeroPage {
		address = uint16(cpu.Memory[cpu.PC+1]) & 0x00FF
	} else if mode == ZeroPageX {
		address = uint16(cpu.Memory[cpu.PC+1]+cpu.X) & 0x00FF
	} else if mode == Absolute {
		address = uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])
	} else if mode == AbsoluteX {
		address = (uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])) + uint16(cpu.X)
		if (address & 0x00FF) < uint16(cpu.X) {
			pageCrossed = true
		}
	} else if mode == AbsoluteY {
		address = (uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])) + uint16(cpu.Y)
		if (address & 0x00FF) < uint16(cpu.Y) {
			pageCrossed = true
		}
	} else if mode == IndexedIndirect {
		intermediateAddress := (uint8(cpu.Memory[cpu.PC+1]) + cpu.X)
		address = uint16(cpu.Memory[intermediateAddress])
	} else if mode == IndirectIndexed {
		intermediateAddress := cpu.Memory[cpu.PC+1]
		address = uint16(cpu.Memory[intermediateAddress]) + uint16(cpu.Y)
		if (address & 0x00FF) < uint16(cpu.Y) {
			pageCrossed = true
		}
	} else if mode == Relative {
		address = cpu.PC + 1
	} else if mode == Implied {
	} else {
		// WIP: garbage value for now, later this should probably be a switch statement
		// and blow up if we don't know the addressing mode
		address = 0x00
	}

	return &InstructionContext{
		Mode:        mode,
		PageCrossed: pageCrossed,
		Address:     address,
	}
}

func NewCPU() *CPU {
	return &CPU{
		PC:     0x00,
		SP:     0xFF,
		A:      0x00,
		X:      0x00,
		Y:      0x00,
		F:      0x00,
		Cycles: 0,
	}
}
