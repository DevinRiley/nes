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

func (cpu *CPU) setZeroFlag() {
	// set zero flag if accumulator is zero
	if cpu.A == 0 {
		cpu.ZFlag = true
	} else {
		cpu.ZFlag = false
	}
}

func (cpu *CPU) setNegativeFlag() {
	// Set negative flag if Accumulator bit 7 is 1
	if cpu.A&byte(0x01<<7) != 0 {
		cpu.NFlag = true
	} else {
		cpu.NFlag = false
	}
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
	} else {
		// WIP: garbage value for now, later this should probably be a switch statement
		// and blow up if we don't know the addressing mode
		address = 0x01
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
		SP:     0x00,
		A:      0x00,
		X:      0x00,
		Y:      0x00,
		F:      0x00,
		Cycles: 0,
	}
}
