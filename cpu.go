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

	if cpu.Debug {
		fmt.Println("Beginning execution...")
		//fmt.Printf("Opcode: %#x, Addressing Mode: %d, Address: %d\n", opcode, context.Mode, context.Address)
		cpu.Print()
	}

	instruction := instructionMap[opcode]
	cpu.PC += instruction.Bytes
	instruction.Exec(cpu, context)
	cpu.Cycles += instruction.Cycles
	if context.PageCrossed && instruction.AddCycleOnPageCross {
		cpu.Cycles += 1
	}

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

func (cpu *CPU) stackPush16(value uint16) {
	cpu.stackPush(byte(value >> 8)) // high byte first
	cpu.stackPush(byte(value))      // then the low byte
}

func (cpu *CPU) stackPop16() uint16 {
	lo := cpu.stackPop()
	hi := cpu.stackPop()
	return uint16(hi)<<8 | uint16(lo)

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
	var pageCrossed = false
	var mode = instructionMap[opcode].AddressingMode

	switch mode {
	case Immediate:
		address = cpu.PC + 1
	case ZeroPage:
		address = uint16(cpu.Memory[cpu.PC+1]) & 0x00FF
	case ZeroPageX:
		address = uint16(cpu.Memory[cpu.PC+1]+cpu.X) & 0x00FF
	case ZeroPageY:
		address = uint16(cpu.Memory[cpu.PC+1]+cpu.Y) & 0x00FF
	case Absolute:
		address = uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])
	case AbsoluteX:
		address = (uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])) + uint16(cpu.X)
		if (address & 0x00FF) < uint16(cpu.X) {
			pageCrossed = true
		}
	case AbsoluteY:
		address = (uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])) + uint16(cpu.Y)
		if (address & 0x00FF) < uint16(cpu.Y) {
			pageCrossed = true
		}
	case IndexedIndirect:
		intermediateAddress := (uint8(cpu.Memory[cpu.PC+1]) + cpu.X)
		address = uint16(cpu.Memory[intermediateAddress])
	case IndirectIndexed:
		intermediateAddress := cpu.Memory[cpu.PC+1]
		address = uint16(cpu.Memory[intermediateAddress]) + uint16(cpu.Y)
		if (address & 0x00FF) < uint16(cpu.Y) {
			pageCrossed = true
		}
	case Relative:
		address = cpu.PC + 1
	case Indirect:
		// NOTE: This addressing mode implements a hardware bug.
		// The operand of the instruction is an intermediate address.
		// This intermediate address contains the low byte of the JMP target address.
		// The high byte of the JMP target address is stored in the following address
		// (intermediate address + 1). However If the intermediate address falls on a
		// memory page boundary (i.e. the first byte is on $xxFF, where xx is any number),
		// it does not correctly look at the next page when reading the high byte of the
		// JMP target address. A concrete example: If the instruction has the operand $10FF,
		// it will read the LSB of the JMP address from $10FF, but will read the MSB of the JMP
		// address from $1000 instead of $1100.
		intermediateLo := uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])
		intermediateHi := (intermediateLo & 0xFF00) | ((intermediateLo + 1) & 0x00FF) // this is the bug
		address = uint16(cpu.Memory[intermediateHi])<<8 | uint16(cpu.Memory[intermediateLo])
	case Implied:
		// this case intentionally left blank :O
	}

	return &InstructionContext{
		PageCrossed:    pageCrossed,
		Address:        address,
		AddressingMode: mode,
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
